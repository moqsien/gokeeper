package kexecutor

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gtree"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/os/gcfg"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
	kapp "github.com/moqsien/gokeeper/kapp"
	"github.com/moqsien/gokeeper/ktype"
	process "github.com/moqsien/processes"
	logger "github.com/moqsien/processes/logger"
)

// Keeper实现本接口的方法在ki_executor.go中
type IKeeper interface {
	Config() *gcfg.Config
	IsMaster() bool
	ListOfAppsToStart() *garray.StrArray
	Mode() ktype.ProcMode
	NewProcess(name string, opts ...process.Option) (*process.ProcessPlus, error)
	Manager() *process.ProcManager
}

/*
Executor 用于保存和运行App；一个Executor可以保存多个App。
在多进程模式下，一个Executor会开启一个新的进程来运行其下的所有App，一个App在新进程中对应一个goroutine。
在单进程模式下，EXecutor只会开启新的goroute来运行行App，所有的goroutine都在一个进程中。
*/
type Executor struct {
	Keeper  IKeeper        // Executor所属的管理者
	Name    string         // 执行器名称
	AppList *gtree.AVLTree // 保存的App列表，key: appName, value: appContainer
}

/*
Executor工厂
*/
func NewExecutor(execName string, k IKeeper) *Executor {
	return &Executor{
		Keeper:  k,
		Name:    execName,
		AppList: gtree.NewAVLTree(gutil.ComparatorString, true),
	}
}

/*
  StopExecutor 停止执行当前Executor；
  会关闭所有正在运行的App。
*/
func (that *Executor) StopExecutor() {
	for _, app := range that.AppList.Map() {
		a := app.(*kapp.AppContainer)
		if a.State == process.Running {
			a.State = process.Stopping
			if e := a.App.Exit(); e != nil {
				logger.Errorf("服务 %s .结束出错，error: %v", a.App.AppName(), e)
			} else {
				logger.Printf("%s 服务 已结束.", a.App.AppName())
			}
			a.State = process.Stopped
			a.StopTime = gtime.Now()
		}
	}
	return
}

// 通过反射生成私有app对象
func (that *Executor) MakeApp(a kapp.IApp) (kapp.IApp, error) {
	var (
		cType  = reflect.TypeOf(a)
		cValue = reflect.ValueOf(a)
	)
	//判断是否是指针类型
	if cType.Kind() != reflect.Ptr {
		return nil, gerror.Newf("生成App: 传入的App对象不是指针类型: %s", cType.String())
	}
	var cTypeElem = cType.Elem()
	//判断是否是struct类型
	if cTypeElem.Kind() != reflect.Struct {
		return nil, gerror.Newf("生成App: 传入的App对象不是struct类型: %s", cType.String())
	}
	//如果结构体没有实现 AppCtx 的方法，或者不是匿名结构体
	iType, ok := cTypeElem.FieldByName("AppBase")
	if !ok || !iType.Anonymous {
		return nil, gerror.Newf("生成App: 传入的App对象未继承 AppBase : %s", cType.String())
	}

	_, found := cType.MethodByName("Execute")
	if !found {
		return nil, gerror.Newf("生成App: 传入的App对象未实现Execute方法")
	}

	_, found = cType.MethodByName("Exit")
	if !found {
		return nil, gerror.Newf("生成App: 传入的App对象未实现Exit方法")
	}

	_, found = cType.MethodByName("Name")
	if !found {
		return nil, gerror.Newf("生成App: 传入的App对象未实现Name方法")
	}
	iValue := cValue.Elem().FieldByName("Executor")
	if iValue.CanSet() {
		iValue.Set(reflect.ValueOf(that))
	}
	iValue = cValue.Elem().FieldByName("Context")
	if iValue.CanSet() {
		iValue.Set(reflect.ValueOf(context.Background()))
	}
	iValue = cValue.Elem().FieldByName("Config")
	if iValue.CanSet() {
		c := &kapp.AppConfig{}
		c.Config = that.Keeper.Config()
		iValue.Set(reflect.ValueOf(c))
	}
	return a, nil
}

// AddApp 添加App到Executor
func (that *Executor) AddApp(a kapp.IApp) error {
	name := a.AppName()
	_, found := that.AppList.Search(name)
	if found {
		return gerror.Newf("App [%s] 已存在", name)
	}
	app, err := that.MakeApp(a)
	if err != nil {
		return err
	}
	that.AppList.Set(app.AppName(), &kapp.AppContainer{
		App:   app,
		State: process.Unknown,
	})
	return nil
}

// RemoveApp 从Executor移除App
func (that *Executor) RemoveApp(name string) {
	value := that.AppList.Remove(name)
	if value == nil {
		return
	}
	app := value.(*kapp.AppContainer)
	if app.State == process.Running {
		err := that.StopApp(name)
		if err != nil {
			logger.Error(err)
		}
	}
}

// StopApp 关闭指定的App
func (that *Executor) StopApp(name string) error {
	a, found := that.AppList.Search(name)
	if !found {
		return fmt.Errorf("未找到[%s]", name)
	}
	ac := a.(*kapp.AppContainer)
	if ac.State == process.Running {
		ac.State = process.Stopping
		err := ac.App.Exit()
		ac.State = process.Stopped
		ac.StopTime = gtime.Now()
		return err
	}
	return nil
}

// StartApp 启动指定的App
func (that *Executor) StartApp(name string) error {
	a, found := that.AppList.Search(name)
	if !found {
		return fmt.Errorf("未找到[%s]", name)
	}
	ac := a.(*kapp.AppContainer)
	if ac.State == process.Starting || ac.State == process.Running {
		return fmt.Errorf("App[%s]正在运行中", name)
	}
	ac.StartTime = gtime.Now()
	ac.State = process.Running
	go func(a1 *kapp.AppContainer) {
		e := a1.App.Execute()
		if e != nil && a1.State != process.Stopping {
			a1.State = process.Stopped
			logger.Warningf("App:[%v] 启动失败: %v", a1.App.AppName(), e)
		}
	}(ac)
	return nil
}

/*
  StartApps
  启动Executor中需要启动的App；
  多进程模式下，本方法在子进程中执行；
  单进程模式下，笨方法在主进程中执行(因为只有一个进程)；
*/
func (that *Executor) StartAllApps() {
	for name, app := range that.AppList.Map() {
		a := app.(*kapp.AppContainer)
		/*
		  如果传入了待启动的AppName，那么需要对比传入的AppName和当前Executor中的AppName，若相同，则启动该App。
		  如果未传入待启动的AppName，则启动当前Executor中的所用App。
		  因此，可以支持启动当前Executor中的一部分App。
		*/
		if that.Keeper.ListOfAppsToStart().Len() > 0 && !that.Keeper.ListOfAppsToStart().ContainsI(a.App.AppName()) {
			that.RemoveApp(gconv.String(name))
			continue
		}
		a.StartTime = gtime.Now()
		a.State = process.Running
		// 尝试启动App
		go func(a1 *kapp.AppContainer) {
			e := a1.App.Execute()
			if e != nil && a1.State != process.Stopping {
				a1.State = process.Stopped
				logger.Warningf("App:[%v] Start Fails: %v", a1.App.AppName(), e)
			}
		}(a)
	}
}

// GetAppCanStart 获取本Executor中需要在子进程中启动的AppNames
func (that *Executor) GetAppNeedToStart() []string {
	/*
	  如果命令行传入了需要启动的AppName，则需要选出当前Executor中的需要启动的AppName；
	  然后将其作为参数传给子进程。子进程根据Executor.StartAllApps方法来启动传过去的App。

	  如果命令行未传入需要启动的AppName，则表示需要启动当前Executor中的所有App。
	*/
	var appList []string
	if that.Keeper.ListOfAppsToStart().Len() > 0 {
		for _, name := range that.AppList.Keys() {
			if that.Keeper.ListOfAppsToStart().ContainsI(gconv.String(name)) {
				appList = append(appList, gconv.String(name))
			}
		}
	} else {
		appList = gconv.Strings(that.AppList.Keys())
	}
	return appList
}

/*
  CreateNewProcess
  启动新的子进程来执行start命令；
  本方法只在主进程中执行；
  一个子进程对应于一个Executor。
*/
func (that *Executor) NewChildProcForStart(configFilePath string) {
	if that.AppList.Size() == 0 {
		return
	}

	/*
	  获取需要在子进程中启动的appNames，
	  如果没有App需要启动，则不会创建新的子进程
	*/
	if appNameList := that.GetAppNeedToStart(); len(appNameList) != 0 {
		// 子进程参数——子进程需要执行的命令：start
		var args = []string{"start"}

		// 子进程参数——运行环境：product, test, dev；默认product
		if len(that.Keeper.Config().GetString("ENV_NAME")) > 0 {
			args = append(args, fmt.Sprintf("--env=%s", that.Keeper.Config().GetString("ENV_NAME")))
		}

		// 子进程参数——配置文件路径
		if len(configFilePath) > 0 {
			args = append(args, fmt.Sprintf("--config=%s", configFilePath))
		}

		// 子进程参数——是否开启调试
		if that.Keeper.Config().GetBool("Debug") {
			args = append(args, "--debug")
		}

		/*
		  子进程参数——子进程中需要启动的Executor；
		  传递本参数之后，明确告诉子进程需要执行哪个Executor，无需遍历所有Executor，提高效率。
		*/
		args = append(args, fmt.Sprintf("--executor=%s", that.Name))

		// 子进程参数——要启动的AppNames
		args = append(args, appNameList...)

		// 	process.ProcExtraFiles(that.server.graceful.getExtraFiles()), // 与获取inheritedEnv的顺序不能错乱
		// 	process.ProcEnvironment(that.server.graceful.inheritedEnv.Map()),

		// 创建新的子进程
		p, e := that.Keeper.NewProcess(that.Name, // 进程名==Executor名称
			process.ProcPath(os.Args[0]),
			process.ProcArgs(args),
			process.ProcEnvVar(ktype.EnvIsChild, "true"),
			process.ProcEnvVar(ktype.EnvIsMaster, "false"), // 子进程的"主进程标记"设置为false，用于区分子进程和主进程
			process.ProcStdoutLog("/dev/stdout", ""),
			process.ProcRedirectStderr(true),
			process.ProcAutoReStart(process.AutoReStartTrue),
			process.ProcStopSignal("SIGQUIT", "SIGTERM"),
			process.ProcStopWaitSecs(int(ktype.MinShutdownTimeout/time.Second)))
		if e != nil {
			logger.Warning(e)
		}

		/*
		  异步开启新的子进程；一个goroutine(在StartProc中实现)对应一个子进程，
		*/
		p.StartProc(true)
	}
}

// Start Executor 执行的调用入口
func (that *Executor) Start(configFilePath string) {
	if that.Keeper.IsMaster() && that.Keeper.Mode() == ktype.MultiProcs {
		/*
		  在多进程模式下，且在主进程中，创建并启动新的子进程
		*/
		that.NewChildProcForStart(configFilePath)
	} else {
		/*
		  1、在多进程模式下，且在子进程中，启动所有需要启动的App；
		  2、在单进程模式下，直接在主进程中启动所有需要启动的App；
		*/
		that.StartAllApps()
	}
}

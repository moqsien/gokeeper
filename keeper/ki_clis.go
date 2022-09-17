package keeper

import (
	"fmt"
	"os"
	"runtime"

	"github.com/gogf/gf"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcfg"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
	gokeeper "github.com/moqsien/gokeeper"
	kexecutor "github.com/moqsien/gokeeper/kexecutor"
	ktype "github.com/moqsien/gokeeper/ktype"
	kutils "github.com/moqsien/gokeeper/kutils"
	logger "github.com/moqsien/processes/logger"
	signals "github.com/moqsien/processes/signals"
	"github.com/spf13/cobra"
)

// 将配置文件中的配置信息写入到logger配置中
func (that *Keeper) InitLogCfg() {
	if !that.KConfig.Available() {
		return
	}
	var m map[string]interface{}
	nodeKey, _ := gutil.MapPossibleItemByKey(that.KConfig.GetMap("."), ktype.ConfigNodeNameLogger)
	if nodeKey == "" {
		nodeKey = ktype.ConfigNodeNameLogger
	}
	m = that.KConfig.GetMap(fmt.Sprintf(`%s.%s`, nodeKey, glog.DefaultName))
	if len(m) == 0 {
		m = that.KConfig.GetMap(nodeKey)
	}
	if len(m) > 0 {
		if err := logger.SetConfigWithMap(m); err != nil {
			panic(err)
		}
	}
}

/*
  以下所有方法实现kcli需要的ICommand接口
*/

// 命令行参数解析和设置
func (that *Keeper) ParseConfig(conf string) {
	that.KConfig = that.GetGFConf(conf)
	that.InitLogCfg() // 配置文件的日志相关配置写入日志配置中
	that.KConfigPath = conf
}

func (that *Keeper) SetRootCommand(c *cobra.Command) {
	that.Command = c
}

func (that *Keeper) SetAppsToOperate(appNames []string) {
	// 设置要操作的App的名称列表
	for _, name := range appNames {
		if gstr.ContainsI(name, ",") {
			appNameArray := gstr.SplitAndTrim(name, ",")
			if len(appNameArray) > 1 {
				that.AppsToOperate.Append(appNameArray...)
			}
			continue
		}
		that.AppsToOperate.Append(gstr.Trim(name))
	}
}

func (that *Keeper) SetCurrentExecutor(execName string) {
	that.CurrentExecutor = execName
}

func (that *Keeper) ParseEnv(env string) {
	if len(env) > 0 {
		// 如果命令行传入了env参数，则使用命令行参数
		_ = genv.Set("ENV_NAME", gstr.ToLower(env))
		_ = that.KConfig.Set("ENV_NAME", gstr.ToLower(env))
	} else if len(that.KConfig.GetString("ENV_NAME")) <= 0 {
		// 如果命令行未传入env参数，且配置文件中也不存在ENV_NAME配置，则先查找环境变量ENV_NAME，并把环境变量中的ENV_NAME赋值给配置文件
		_ = that.KConfig.Set("ENV_NAME", gstr.ToLower(genv.Get("ENV_NAME", "product")))
	}
}

func (that *Keeper) ParseDebug(debug bool) {
	if debug {
		// 如果启动命令行强制设置了debug参数，则优先级最高
		_ = genv.Set("DEBUG", "true")
		_ = that.KConfig.Set("Debug", true)
	} else {
		// 默认设置为false
		_ = that.KConfig.Set("Debug", that.KConfig.GetBool("Debug", genv.GetVar("DEBUG", false).Bool()))
	}
}

func (that *Keeper) ParseDaemon(daemon bool) {
	_ = that.KConfig.Set("Daemon", daemon)
}

func (that *Keeper) ParsePidFilePath(pidFilePath string) {
	// 主进程pid文件保存路径
	if len(pidFilePath) > 0 {
		that.PidFilePath = pidFilePath
		return
	}
	that.PidFilePath = gfile.TempDir(fmt.Sprintf("%s.pid", that.KeeperName))
}

/*
  日志初始化
  通过参数设置日志级别
  日志级别通过环境默认分三个类型，开发环境，测试环境，生产环境
  开发环境: 日志级别为 DEVELOP,标准输出打开
  测试环境：日志级别为 INFO,除了debug日志，都会被打印，标准输出关闭
  生产环境: 日志级别为 PRODUCT，会打印 WARN,ERRO,CRIT三个级别的日志，标准输出为关闭
  Debug开关会无视以上设置，强制把日志级别设置为ALL，并且打开标准输出。
*/
func (that *Keeper) InitLogSetting(config *gcfg.Config) error {
	loggerCfg := config.GetJson("logger")
	env := config.GetString("ENV_NAME")
	level := loggerCfg.GetString("Level")
	logger.SetDebug(false)
	logger.SetStdoutPrint(false)
	//如果配置文件中的日志配置不存在，则判断环境变量，通过不同的环境变量，给与不同的日志级别
	if len(level) <= 0 {
		if env == "dev" || env == "develop" {
			level = "DEVELOP"
		} else if env == "test" {
			level = "INFO"
		} else {
			level = "PRODUCT"
		}
	}

	setConfig := g.Map{"level": level}

	if env == "dev" || env == "develop" {
		setConfig["stdout"] = true
		logger.SetDebug(true)
	}
	logPath := loggerCfg.GetString("Path")
	if len(logPath) > 0 {
		setConfig["path"] = logPath
	} else {
		logger.SetDebug(true)
	}

	// 如果开启debug模式，则无视其他设置
	if config.GetBool("Debug", false) {
		setConfig["level"] = "ALL"
		setConfig["stdout"] = true
		logger.SetDebug(true)
	}
	return logger.SetConfigWithMap(setConfig)
}

/*
  CheckKeeperStart 检查keeper是否已经启动过；
  如果主进程的pid文件名存在，则说明keeper已经启动过。
*/
func (that *Keeper) CheckKeeperForStart() {
	pidFile := that.PidFilePath
	var keeperPid = 0
	if gfile.IsFile(pidFile) {
		keeperPid = gconv.Int(gstr.Trim(gfile.GetContents(pidFile)))
	}
	if keeperPid == 0 {
		return
	}
	if signals.CheckPidExist(keeperPid) {
		logger.Fatalf("Keeper [%d] is already running.", keeperPid)
	}
	return
}

/*
  RunExecutors 根据命令行参数执行ExecutorList中的Executor
*/
func (that *Keeper) RunExecutors() {
	if that.IsMutilProcModeAndInMaster() {
		// TODO: 多进程模式下，且在主进程中，创建新的子进程
		// err := that.inheritListenerList()
		that.ExecutorList.Iterator(func(_ interface{}, v interface{}) bool {
			ke := v.(*kexecutor.Executor)
			// 会将对应的Executor名称和需要启动的App传给子进程
			ke.NewChildProcForStart(that.KConfigPath)
			return true
		})
	} else if that.ProcMode == ktype.MultiProcs && !that.IsMaster() {
		// 多进程模式下，且在子进程中，执行对应的Executor
		if exec, existed := that.ExecutorList.Search(that.CurrentExecutor); existed {
			executor, ok := exec.(kexecutor.Executor)
			if ok {
				executor.StartAllApps()
			}
		}
	} else {
		// 单进程模式下，直接在主进程中启动所有App
		that.ExecutorList.Iterator(func(_ interface{}, v interface{}) bool {
			ke := v.(*kexecutor.Executor)
			ke.StartAllApps()
			return true
		})
	}
}

// RunKeeper keeper的start命令的执行入口
func (that *Keeper) RunKeeper() {
	//判断是否是守护进程运行
	if e := kutils.Demonize(that.KConfig); e != nil {
		logger.Fatalf("error:%v", e)
	}
	//记录启动时间
	that.StartTime = gtime.Now()

	// TODO: 初始化平滑重启的钩子函数
	// that.InitGraceful()

	/*
	  执行业务入口函数，在Keeper.Setup方法中设置；
	  keeper_default.Setup方法为默认keeper设置startfunc；
	*/
	if that.StartFunction == nil {
		logger.Error("StartFunction has not been set!")
		return
	}
	that.StartFunction(that)

	// TODO: 设置优雅退出时候需要做的工作
	// that.Graceful.SetShutdown(15*time.Second, that.FirstStop, that.BeforeExiting)

	// 执行ExecutorList中的Executor
	that.RunExecutors()

	if that.ProcMode == ktype.SingleProc || that.IsMutilProcModeAndInMaster() {
		that.PutMasterPidInFile()
	}

	logger.Printf("%d: 服务已经初始化完成, %d 个协程被创建.", os.Getpid(), runtime.NumGoroutine())
}

// StopKeeper keeper的stop、reload、quit命令的执行入口
func (that *Keeper) StopKeeper(sig string) {
	pidFile := that.PidFilePath
	var keeperPid = 0
	if gfile.IsFile(pidFile) {
		keeperPid = gconv.Int(gstr.Trim(gfile.GetContents(pidFile)))
	}
	if keeperPid == 0 {
		logger.Println("Keeper is not running.")
		os.Exit(0)
	}

	var sigNo string
	switch sig {
	case "stop":
		sigNo = "SIGTERM"
	case "reload":
		sigNo = "SIGUSR2"
	case "quit":
		sigNo = "SIGQUIT"
	default:
		logger.Printf("signal cmd `%s' not found", sig)
		os.Exit(0)
	}
	err := signals.KillPid(keeperPid, signals.ToSignal(sigNo), false)
	if err != nil {
		logger.Printf("error:%v", err)
	}
	os.Exit(0)
}

/*
  Version相关信息
*/
var (
	BuildVersion     = "No Version Info"
	BuildGoVersion   = "No Version Info"
	BuildGitCommitId = "No Commit Info"
	BuildTime        = "No Time Info"
	Authors          = "No Authors Info"
	Logo             = `
	_                                  
	| | __ ___   ___  _ __    ___  _ __ 
	| |/ // _ \ / _ \| '_ \  / _ \| '__|
	|   <|  __/|  __/| |_) ||  __/| |   
	|_|\_\\___| \___|| .__/  \___||_|   
					 |_|                   
`
)

// Version 显示版本信息
func (that *Keeper) Version() {
	fmt.Print(gstr.TrimLeftStr(Logo, "\n"))
	fmt.Printf("Version:         %s\n", BuildVersion)
	fmt.Printf("Go Version:      %s\n", runtime.Version())
	fmt.Printf("GoKeeper Author: %s\n", gokeeper.Author)
	fmt.Printf("GoKeeper Version:%s\n", gokeeper.Version)
	fmt.Printf("GF Version:      %s\n", gf.VERSION)
	fmt.Printf("Git Commit:      %s\n", BuildGitCommitId)
	fmt.Printf("Build Time:      %s\n", BuildTime)
	fmt.Printf("Authors:         %s\n", Authors)
	fmt.Printf("Install Path:    %s\n", gfile.SelfPath())
}

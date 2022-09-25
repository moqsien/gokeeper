package keeper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	kexecutor "github.com/moqsien/gokeeper/kexecutor"
	ktype "github.com/moqsien/gokeeper/ktype"
	kutils "github.com/moqsien/gokeeper/kutils"
	goktrl "github.com/moqsien/goktrl"
	logger "github.com/moqsien/processes/logger"
)

func (that *Keeper) KCtrlCheckExecutor(c *goktrl.Context) bool {
	eName := c.Parser.GetOpt("executor")
	_, ok := that.Manager.Search(eName)
	if !ok {
		fmt.Printf("Executor: [%s] does not exist", eName)
	}
	return ok
}

/* TODO:
交互式shell初始化
*/

func (that *Keeper) kCtrlVersion() {
	versionFun := func(k *goktrl.Context) {
		that.Version()
	}

	that.KCtrl.AddKtrlCommand(&goktrl.KCommand{
		Name:       "version",
		Help:       "show keeper version info.",
		Func:       versionFun,
		SocketName: that.KCtrlSocket,
	})
}

func (that *Keeper) kCtrlInfo() {
	type Data struct {
		Keeper     string `order:"1"`
		ProcMode   string `order:"2"`
		Executor   string `order:"3"`
		Pid        int    `order:"4"`
		Apps       string `order:"6"`
		AppsRunnig string `order:"5"`
	}

	var Result = []*Data{} // 客户端和服务端在不同进程中，此处无影响

	info := func(c *goktrl.Context) {
		result, err := c.GetResult()
		if err != nil {
			logger.Error(err)
			return
		}
		err = json.Unmarshal(result, &Result)
		if err != nil {
			logger.Error(err)
			return
		}
		c.Table.AddRowsByListObject(Result)
	}

	handler := func(c *goktrl.Context) {
		that.Manager.Iterator(func(_ string, v interface{}) bool {
			executor := v.(*kexecutor.Executor)
			Result = append(Result, &Data{
				Keeper:     that.KeeperName,
				ProcMode:   that.ProcMode.String(),
				Executor:   executor.Name,
				Pid:        executor.Pid,
				AppsRunnig: executor.AppsRunning.String(),
				Apps:       kutils.ConvertSliceToString(executor.AppList.Keys()),
			})
			return true
		})

		if content, err := json.Marshal(Result); err == nil {
			c.String(http.StatusOK, string(content))
		} else {
			logger.Error(err)
		}
	}

	that.KCtrl.AddKtrlCommand(&goktrl.KCommand{
		Name:        "info",
		Help:        "show keeper info",
		Func:        info,
		ShowTable:   true,
		KtrlHandler: handler,
		SocketName:  that.KCtrlSocket,
	})
}

// KtrlStartExecutor 启动一个Executor，可以指定启动一部分app
func (that *Keeper) KtrlStartExecutor() {
	type OptsStartExecutor struct {
		Executor string `alias:"e" required:"true" descr:"executor from keeper."`
	}
	start := func(c *goktrl.Context) {
		// 单进程不能启动新Executor
		if that.IsSingleProcMode() {
			fmt.Println("Single-process mode!")
			return
		}
		// 必须传入存在的Executor名称
		if _, found := that.Manager.Search(c.Parser.GetOpt("executor")); !found {
			fmt.Printf("Executor: [%s] does not exist", c.Parser.GetOpt("executor"))
			return
		}
		result, err := c.GetResult()
		if err != nil {
			logger.Error(err)
			return
		}
		fmt.Println(string(result))
	}

	handler := func(c *goktrl.Context) {
		// TODO: appNames
		if eName := c.Query("executor"); len(eName) > 0 {
			c.String(http.StatusOK, that.StartExecutor(eName))
		}
	}

	that.KCtrl.AddKtrlCommand(&goktrl.KCommand{
		Name:            "startexc",
		Help:            "start an executor.",
		Func:            start,
		Opts:            &OptsStartExecutor{},
		KtrlHandler:     handler,
		SocketName:      that.KCtrlSocket,
		ArgsDescription: "apps to start.",
	})
}

// KtrlStartApp 启动App，注意必须指定一个已经启动的Executor
func (that *Keeper) KtrlStartApps() {
	type OptsStartApps struct {
		Executor string `alias:"e" required:"true" descr:"executor from keeper."`
	}
	start := func(c *goktrl.Context) {
		var (
			result []byte
			err    error
		)
		if that.ProcMode == ktype.MultiProcs {
			// 多进程模式下，必须传入存一个在的Executor名称
			result, err = c.GetResult(c.Parser.GetOpt("executor"))
		} else {
			// 单进程模式下启动apps, Executor任意传
			result, err = c.GetResult()
		}
		if err != nil {
			if strings.Contains(err.Error(), "connect: no such file or directory") {
				fmt.Printf("Executor: %v is not running!\n", c.Parser.GetOpt("executor"))
				return
			}
			fmt.Println(err)
			return
		}
		fmt.Println(string(result))
	}
	handler := func(c *goktrl.Context) {

	}

	that.KCtrl.AddKtrlCommand(&goktrl.KCommand{
		Name:            "startapp",
		Help:            "start apps.",
		Func:            start,
		Opts:            &OptsStartApps{},
		KtrlHandler:     handler,
		SocketName:      that.KCtrlSocket,
		ArgsRequired:    true,
		ArgsDescription: "apps to start.",
	})
}

// KtrlStopExecutor 停止一个Executor
func (that *Keeper) KtrlStopExecutor() {
	type OptsStopExecutor struct {
		Executor string `alias:"e" required:"true" descr:"executor from keeper."`
	}
	stop := func(k *goktrl.Context) {

	}
	handler := func(c *goktrl.Context) {

	}

	that.KCtrl.AddKtrlCommand(&goktrl.KCommand{
		Name:        "stopexc",
		Help:        "stop an Executor.",
		Func:        stop,
		Opts:        &OptsStopExecutor{},
		KtrlHandler: handler,
		SocketName:  that.KCtrlSocket,
	})
}

func (that *Keeper) KtrlStopApps() {
	type OptsStopApps struct {
		Executor string `alias:"e" required:"true" descr:"executor from keeper."`
	}
	stop := func(k *goktrl.Context) {

	}
	handler := func(c *goktrl.Context) {

	}

	that.KCtrl.AddKtrlCommand(&goktrl.KCommand{
		Name:        "stopapp",
		Help:        "stop apps.",
		Func:        stop,
		Opts:        &OptsStopApps{},
		KtrlHandler: handler,
		SocketName:  that.KCtrlSocket,
	})
}

func (that *Keeper) KtrlReload() {
	reload := func(k *goktrl.Context) {}
	handler := func(c *goktrl.Context) {}

	that.KCtrl.AddKtrlCommand(&goktrl.KCommand{
		Name:        "reload",
		Help:        "reload apps or an executor.",
		Func:        reload,
		KtrlHandler: handler,
		SocketName:  that.KCtrlSocket,
	})
}

func (that *Keeper) KtrlDebug() {
	debug := func(k *goktrl.Context) {}
	handler := func(c *goktrl.Context) {}

	that.KCtrl.AddKtrlCommand(&goktrl.KCommand{
		Name:        "debug",
		Help:        "set debug mode",
		Func:        debug,
		KtrlHandler: handler,
		SocketName:  that.KCtrlSocket,
	})
}

func (that *Keeper) KtrlLog() {

}

func (that *Keeper) InitKtrl() {
	if that.KeeperIsMaster {
		that.KCtrlSocket = that.KeeperName
	} else if that.CurrentExecutor != "" {
		that.KCtrlSocket = that.CurrentExecutor
	}
	if !that.IsCtrlInitiated {
		that.kCtrlVersion()
		that.kCtrlInfo()
		that.KtrlStartExecutor()
		that.KtrlStartApps()
		that.KtrlStopExecutor()
		that.KtrlStartApps()
		that.KtrlReload()
		that.KtrlDebug()
		that.KtrlLog()
	}
	that.IsCtrlInitiated = true // KCtrl标记为已初始化
}

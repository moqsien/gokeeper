package keeper

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/frame/g"
	kexecutor "github.com/moqsien/gokeeper/kexecutor"
	"github.com/moqsien/gokeeper/ktype"
	kutils "github.com/moqsien/gokeeper/kutils"
	goktrl "github.com/moqsien/goktrl"
	logger "github.com/moqsien/processes/logger"
)

func (that *Keeper) KCtrlCheckExecutor(k *goktrl.KtrlContext) bool {
	eName := k.Parser.GetOpt("executor")
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
	versionFun := func(k *goktrl.KtrlContext) {
		that.Version()
	}

	that.KCtrl.AddKtrlCommand(&goktrl.KCommand{
		Name:        "version",
		Help:        "show keeper version info",
		Func:        versionFun,
		KtrlHandler: func(c *gin.Context) {},
		SocketName:  that.KCtrlSocket,
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

	info := func(k *goktrl.KtrlContext) {
		result, err := k.GetResult()
		if err != nil {
			logger.Error(err)
			return
		}
		err = json.Unmarshal(result, &Result)
		if err != nil {
			logger.Error(err)
			return
		}
		k.Table.AddRowsByListObject(Result)
	}

	handler := func(c *gin.Context) {
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
	start := func(k *goktrl.KtrlContext) {
		// 单进程不能启动新Executor
		if that.IsSingleProcMode() {
			fmt.Println("Single-process mode!")
			return
		}
		// 必须传入存在的Executor名称
		if ok := that.KCtrlCheckExecutor(k); !ok {
			return
		}
		result, err := k.GetResult()
		if err != nil {
			logger.Error(err)
			return
		}
		fmt.Println(string(result))
	}

	handler := func(c *gin.Context) {
		if eName := c.Query("executor"); len(eName) > 0 {
			c.String(http.StatusOK, that.StartExecutor(eName))
		}
	}

	that.KCtrl.AddKtrlCommand(&goktrl.KCommand{
		Name: "startexc",
		Help: "【Start a not running Executor】 Usage: startexc -e=<executor> [-a=<app,app1>]",
		Func: start,
		Opts: &g.MapStrBool{
			"executor,e": true,
			"apps,app,a": true, // 可选，如果传了，则只启动指定的App，没传则启动所有App
		},
		KtrlHandler: handler,
		SocketName:  that.KCtrlSocket,
	})
}

// KtrlStartApp 启动App，注意必须指定一个已经启动的Executor
func (that *Keeper) KtrlStartApps() {
	start := func(k *goktrl.KtrlContext) {
		if that.ProcMode == ktype.MultiProcs {
			// 多进程模式下，必须传入存一个在的Executor名称
			if ok := that.KCtrlCheckExecutor(k); !ok {
				return
			}
		} else {
			// 单进程模式下启动apps，无需传入Executor名称
		}
	}
	handler := func(c *gin.Context) {

	}

	that.KCtrl.AddKtrlCommand(&goktrl.KCommand{
		Name: "startapp",
		Help: "【Start apps】 Usage: startapp -e=<executor> -a=<app,app1>",
		Func: start,
		Opts: &g.MapStrBool{
			"executor,e": true, // 必须指定app所属的executor
			"apps,app,a": true,
		},
		KtrlHandler: handler,
		SocketName:  that.KCtrlSocket,
	})
}

// KtrlStopExecutor 停止一个Executor
func (that *Keeper) KtrlStopExecutor() {
	stop := func(k *goktrl.KtrlContext) {

	}
	handler := func(c *gin.Context) {

	}

	that.KCtrl.AddKtrlCommand(&goktrl.KCommand{
		Name: "stopexc",
		Help: "【Stop an Executor】 Usage: stopexc -e=<executor>",
		Func: stop,
		Opts: &g.MapStrBool{
			"executor,e": true,
		},
		KtrlHandler: handler,
		SocketName:  that.KCtrlSocket,
	})
}

func (that *Keeper) KtrlStopApps() {
	stop := func(k *goktrl.KtrlContext) {

	}
	handler := func(c *gin.Context) {

	}

	that.KCtrl.AddKtrlCommand(&goktrl.KCommand{
		Name: "stopapp",
		Help: "【Stop apps】 Usage: stopapp -e=<executor> -a=<app,app1>",
		Func: stop,
		Opts: &g.MapStrBool{
			"executor,e": true,
			"apps,app,a": true,
		},
		KtrlHandler: handler,
		SocketName:  that.KCtrlSocket,
	})
}

func (that *Keeper) KtrlReload() {
	reload := func(k *goktrl.KtrlContext) {}
	handler := func(c *gin.Context) {}

	that.KCtrl.AddKtrlCommand(&goktrl.KCommand{
		Name:        "reload",
		Help:        "reload apps or an executor",
		Func:        reload,
		KtrlHandler: handler,
		SocketName:  that.KCtrlSocket,
	})
}

func (that *Keeper) KtrlDebug() {
	debug := func(k *goktrl.KtrlContext) {}
	handler := func(c *gin.Context) {}

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

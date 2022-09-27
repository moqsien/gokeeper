package keeper

import (
	"fmt"
	"strings"

	kexecutor "github.com/moqsien/gokeeper/kexecutor"
	kutils "github.com/moqsien/gokeeper/kutils"
	goktrl "github.com/moqsien/goktrl"
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
	that.KCtrl.AddKtrlCommand(&goktrl.KCommand{
		Name: "version",
		Help: "show keeper version info.",
		Func: func(k *goktrl.Context) {
			that.Version()
		},
		SocketName: that.KCtrlSocket,
		Auto:       false, // 不需要发请求
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

	that.KCtrl.AddKtrlCommand(&goktrl.KCommand{
		Name: "info",
		Help: "show keeper info",
		KtrlHandler: func(c *goktrl.Context) {
			that.Manager.Iterator(func(_ string, v interface{}) bool {
				executor := v.(*kexecutor.Executor)
				Result = append(Result, &Data{
					Keeper:     that.KeeperName,
					ProcMode:   that.ProcMode.String(),
					Executor:   executor.Name,
					Pid:        executor.Pid,
					Apps:       kutils.SliceToString(executor.AppList.Keys()),
					AppsRunnig: kutils.SliceToString(executor.AppsRunning.Keys()),
				})
				return true
			})
			c.Send(Result)
		},
		Auto:        true,
		ShowTable:   true,
		TableObject: &Result,
		SocketName:  that.KCtrlSocket,
	})
}

// KtrlStartExecutor 启动一个Executor，可以指定启动一部分app
func (that *Keeper) KtrlStartExecutor() {
	type OptsStartExecutor struct {
		Executor string `alias:"e" required:"true" descr:"executor from keeper."`
	}
	that.KCtrl.AddKtrlCommand(&goktrl.KCommand{
		Name: "starte",
		Help: "start an executor.",
		KtrlHandler: func(c *goktrl.Context) {
			opt := c.Options.(*OptsStartExecutor)
			c.Send(that.StartExecutor(opt.Executor, c.Args...))
		},
		Opts:            &OptsStartExecutor{},
		ArgsDescription: "apps to start.",
		Auto:            true,
		SocketName:      that.KCtrlSocket,
	})
}

// KtrlStartApp 启动App
func (that *Keeper) KtrlStartApps() {

	type OptsStartApps struct {
		Executor string `alias:"e" required:"true" descr:"executor from keeper."`
	}
	handler := func(c *goktrl.Context) {
		opt := c.Options.(*OptsStartApps)
		exec, found := that.Manager.Search(opt.Executor)
		if !found {
			c.Send(fmt.Sprintf("Executor: %s is not found!", opt.Executor))
		} else {
			if that.IsMutilProcModeAndInMaster() {
				ex := exec.(*kexecutor.Executor)
				if ex.ProcessPlus != nil && ex.IsRunning() {
					result, _ := c.GetResult(opt.Executor) // 转发给子进程，由子进程运行app
					for _, v := range strings.Split(string(result), ",") {
						ex.AppsRunning.Set(v, struct{}{})
					}
					c.Send(fmt.Sprintf("Apps: [%s] started running.", string(result)))
				} else {
					ex.ProcessPlus = nil
					c.Send(that.StartExecutor(opt.Executor, c.Args...)) // 启动新进程来运行app
				}
			} else {
				ex := exec.(*kexecutor.Executor)
				started := []string{}
				var err error
				for _, name := range c.Args {
					err = ex.StartApp(name)
					if err == nil {
						started = append(started, name)
					}
				}
				if that.IsSingleProcMode() {
					c.Send(fmt.Sprintf("Apps: [%s] started running.", kutils.SliceToString(started)))
				} else {
					c.Send(kutils.SliceToString(started))
				}
			}
		}
	}

	that.KCtrl.AddKtrlCommand(&goktrl.KCommand{
		Name:            "starta",
		Help:            "start apps.",
		KtrlHandler:     handler,
		Auto:            true,
		ArgsRequired:    true,
		ArgsDescription: "apps to start.",
		Opts:            &OptsStartApps{},
		SocketName:      that.KCtrlSocket,
	})
}

// KtrlStopExecutor 停止一个Executor
func (that *Keeper) KtrlStopExecutor() {
	type OptsStopExecutor struct {
		Executor string `alias:"e" required:"true" descr:"executor from keeper."`
	}
	handler := func(c *goktrl.Context) {
		opt := c.Options.(*OptsStopExecutor)
		if exec, found := that.Manager.Search(opt.Executor); !found {
			c.Send(fmt.Sprintf("Executor: %s is not found!", opt.Executor))
		} else {
			ex := exec.(*kexecutor.Executor)
			ex.StopProc(true)
			ex.AppsRunning.Clear()
			ex.Pid = 0
			that.ExecutorsRunning.Remove(ex.Name)
			c.Send(fmt.Sprintf("Executor: %s stopped!", opt.Executor))
		}
	}

	that.KCtrl.AddKtrlCommand(&goktrl.KCommand{
		Name:        "stope",
		Help:        "stop an Executor.",
		Opts:        &OptsStopExecutor{},
		KtrlHandler: handler,
		SocketName:  that.KCtrlSocket,
		Auto:        true,
	})
}

func (that *Keeper) KtrlStopApps() {
	type OptsStopApps struct {
		Executor string `alias:"e" required:"true" descr:"executor from keeper."`
	}
	handler := func(c *goktrl.Context) {
		opt := c.Options.(*OptsStopApps)
		exec, found := that.Manager.Search(opt.Executor)
		if !found {
			c.Send(fmt.Sprintf("Executor: %s is not found!", opt.Executor))
		} else {
			if that.IsMutilProcModeAndInMaster() {
				ex := exec.(*kexecutor.Executor)
				if ex.ProcessPlus != nil && ex.IsRunning() {
					result, _ := c.GetResult(opt.Executor) // 转发给子进程，由子进程运行app
					for _, v := range strings.Split(string(result), ",") {
						ex.AppsRunning.Remove(v)
					}
					c.Send(fmt.Sprintf("Apps: [%s] stopped running.", string(result)))
				} else {
					c.Send(fmt.Sprintf("Executor: %s is not running!", opt.Executor))
				}
			} else {
				ex := exec.(*kexecutor.Executor)
				stopped := []string{}
				var err error
				for _, name := range c.Args {
					err = ex.StopApp(name)
					if err == nil {
						stopped = append(stopped, name)
					}
				}
				if that.IsSingleProcMode() {
					c.Send(fmt.Sprintf("Apps: [%s] stopped running.", kutils.SliceToString(stopped)))
				} else {
					c.Send(kutils.SliceToString(stopped))
				}
			}
		}
	}

	that.KCtrl.AddKtrlCommand(&goktrl.KCommand{
		Name:            "stopa",
		Help:            "stop apps.",
		Opts:            &OptsStopApps{},
		KtrlHandler:     handler,
		SocketName:      that.KCtrlSocket,
		ArgsRequired:    true,
		ArgsDescription: "apps to stop.",
		Auto:            true,
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

package keeper

import (
	"fmt"

	"github.com/moqsien/gokeeper/kexecutor"
)

/*
  交互式shell命令服务端调用的相关的Keeper方法
*/

// StartExecutor 交互式shell开启Executor，只在主进程中执行
func (that *Keeper) StartExecutor(execName string, appNames ...string) (r string) {
	if that.IsMutilProcModeAndInMaster() {
		if _, ok := that.Manager.Search(execName); !ok {
			r = fmt.Sprintf("Executor: [%s] is already running", execName)
		} else {
			executor, _ := that.Manager.Search(execName)
			if that.IsMutilProcModeAndInMaster() {
				that.SetAppsToOperate(appNames)
				executor.(*kexecutor.Executor).NewChildProcForStart(that.KConfigPath)
			}
			r = fmt.Sprintf("Executor: [%s] started!", execName)
		}
	}
	return
}

// StartApps 交互式shell启动Apps，只在子进程中或者单进程模式下执行
func (that *Keeper) StartApps(execName string, appNames ...string) (r string) {
	if !that.IsMaster() && !that.IsSingleProcMode() {

	} else if that.IsSingleProcMode() {

	}
	return
}

func (that *Keeper) StopExecutor(execName string) (r string) {
	return
}

func (that *Keeper) StopApps(appName string) (r string) {
	return
}

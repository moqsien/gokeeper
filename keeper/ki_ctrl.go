package keeper

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/moqsien/gokeeper/kexecutor"
	"github.com/moqsien/gokeeper/kutils"
	"github.com/moqsien/goktrl"
	logger "github.com/moqsien/processes/logger"
)

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
		KtrlPath:    "/ctrl/version",
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

	infoFunc := func(k *goktrl.KtrlContext) {
		result, err := k.GetResult()
		if err != nil {
			logger.Error(err)
			return
		}
		err = json.Unmarshal([]byte(result), &Result)
		if err != nil {
			logger.Error(err)
			return
		}
		k.Table.AddRowsByListObject(Result)
	}

	infoHandler := func(c *gin.Context) {
		that.ExecutorList.Iterator(func(_ interface{}, v interface{}) bool {
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
		Func:        infoFunc,
		KtrlPath:    "/ctrl/info",
		ShowTable:   true,
		KtrlHandler: infoHandler,
		SocketName:  that.KCtrlSocket,
	})
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
	}
	that.IsCtrlInitiated = true // KCtrl标记为已初始化
}

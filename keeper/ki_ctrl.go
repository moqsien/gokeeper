package keeper

import (
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/frame/g"
	"github.com/moqsien/goktrl"
)

func (that *Keeper) kCtrlVersion() {

}

func (that *Keeper) kCtrlInfo() {
	// info命令
	infoFunc := func(k *goktrl.KtrlContext) {

	}

	infoHandler := func(c *gin.Context) {

	}
	that.KCtrl.AddKtrlCommand(&goktrl.KCommand{
		Name:        "info",
		Help:        "show keeper info",
		Func:        infoFunc,
		Opts:        &g.MapStrBool{},
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

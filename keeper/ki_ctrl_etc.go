package keeper

import (
	"net/http"

	"github.com/abiosoft/ishell/v2"
	"github.com/gin-gonic/gin"
	"github.com/moqsien/gokeeper/kctrl"
)

/*
  CtrlKeeperName 方法实现ICtrlKeeper接口
*/
func (that *Keeper) CtrlKeeperName() string {
	return that.KeeperName
}

type CtrlCommand struct {
	CmdName string
	UrlPath string
	Help    string
}

var (
	CmdVersion *CtrlCommand = &CtrlCommand{
		CmdName: "version",
		UrlPath: "/ctrl/version/",
		Help:    "keeper version",
	}
	CmdInfo *CtrlCommand = &CtrlCommand{
		CmdName: "info",
		UrlPath: "/ctrl/info/",
		Help:    "keeper info",
	}
)

// =========================================
/*
  添加交互式Shell服务端路由
*/
func (that *Keeper) AddCtrlRoutes() {
	that.Controller.GET(CmdInfo.UrlPath, func(c *gin.Context) {
		// TODO: info命令
		c.String(http.StatusOK, "Info")
	})
}

/*
  启动交互式Shell服务端
*/
func (that *Keeper) StartCtrlAsServer() {
	if that.Controller != nil && that.Controller.KcMode == kctrl.CtrlUnknown {
		that.Controller.InitServer()      // 初始化为Shell服务端
		that.AddCtrlRoutes()              // 为Shell服务端添加路由
		that.Controller.CtrlServerStart() // 启动Shell服务端
	}
}

// =========================================
/*
  添加交互式Shell命令
*/
func (that *Keeper) AddCtrlCommand() {
	that.Controller.AddCmd(&ishell.Cmd{
		Name: CmdVersion.CmdName,
		Help: CmdVersion.Help,
		Func: func(c *ishell.Context) {
			that.Version()
		},
	})

	that.Controller.AddCmd(&ishell.Cmd{
		Name: CmdInfo.CmdName,
		Help: CmdInfo.Help,
		Func: func(c *ishell.Context) {
			params := kctrl.ParamsContainer{}
			result, err := that.Controller.CtrlGetStr(CmdInfo.UrlPath, params)
			if err != nil {
				c.Println(err)
				return
			}
			c.Println(result)
		},
	})

}

/*
  启动交互式Shell
*/
func (that *Keeper) StartCtrlAsClient() {
	if that.Controller != nil && that.Controller.KcMode == kctrl.CtrlUnknown {
		that.Controller.InitClient()      // 初始化为Shell客户端
		that.AddCtrlCommand()             // 添加Shell命令
		that.Controller.CtrlClientStart() // 启动交互式Shell
	}
}

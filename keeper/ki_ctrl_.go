package keeper

import (
	"net/http"
	"os"

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

func (that *Keeper) CtrlCurrentExecutor() string {
	return that.CurrentExecutor
}

func (that *Keeper) CtrlIsMaster() bool {
	return that.KeeperIsMaster
}

/*
  以下和Keeper.Controller相关，用于创建交互式Shell命令
*/
type CtrlCommand struct {
	CmdName string
	UrlPath string
	Help    string
}

var (
	CmdVersion *CtrlCommand = &CtrlCommand{
		CmdName: "version",
		UrlPath: "/ctrl/version/", // version命令无需请求，UrlPath可以不配置
		Help:    "keeper version",
	}
	CmdInfo *CtrlCommand = &CtrlCommand{
		CmdName: "info",
		UrlPath: "/ctrl/info/",
		Help:    "keeper info",
	}
	CmdStart *CtrlCommand = &CtrlCommand{
		CmdName: "start",
		UrlPath: "/ctrl/start/",
		Help:    "start app",
	}
	CmdStop *CtrlCommand = &CtrlCommand{
		CmdName: "stop",
		UrlPath: "/ctrl/stop/",
		Help:    "stop app",
	}
	CmdReload *CtrlCommand = &CtrlCommand{
		CmdName: "reload",
		UrlPath: "/ctrl/reload/",
		Help:    "reload app",
	}
	CmdDebug *CtrlCommand = &CtrlCommand{
		CmdName: "debug",
		UrlPath: "/ctrl/debug/",
		Help:    "debug settings",
	}
	CmdLog *CtrlCommand = &CtrlCommand{
		CmdName: "log",
		UrlPath: "/ctrl/log/",
		Help:    "print log",
	}
)

// =========================================
/*
  添加交互式Shell服务端路由
*/
func (that *Keeper) AddCtrlRoutes() {
	// TODO: info命令
	that.Controller.GET(CmdInfo.UrlPath, func(c *gin.Context) {
		c.String(http.StatusOK, "Info")
	})

	// TODO: start命令
	that.Controller.GET(CmdStart.UrlPath, func(c *gin.Context) {

	})

	// TODO: stop命令
	that.Controller.GET(CmdStop.UrlPath, func(c *gin.Context) {

	})

	// TODO: reload命令
	that.Controller.GET(CmdReload.UrlPath, func(c *gin.Context) {

	})

	// TODO: debug命令
	that.Controller.GET(CmdDebug.UrlPath, func(c *gin.Context) {

	})

	// TODO: log命令
	that.Controller.GET(CmdLog.UrlPath, func(c *gin.Context) {

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

	// TODO: info命令
	that.Controller.AddCmd(&ishell.Cmd{
		Name: CmdInfo.CmdName,
		Help: CmdInfo.Help,
		Func: func(c *ishell.Context) {
			result, err := that.Controller.CtrlGetStr(CmdInfo.UrlPath, nil)
			if err != nil {
				c.Println(err)
				return
			}
			c.Println(result)
		},
	})

	// TODO: start命令
	that.Controller.AddCmd(&ishell.Cmd{
		Name: CmdStart.CmdName,
		Help: CmdStart.Help,
		Func: func(c *ishell.Context) {
			os.Args = c.Args
		},
	})

	// TODO: stop命令
	that.Controller.AddCmd(&ishell.Cmd{
		Name: CmdStop.CmdName,
		Help: CmdStop.Help,
		Func: func(c *ishell.Context) {
			os.Args = c.Args
		},
	})

	// TODO: reload命令
	that.Controller.AddCmd(&ishell.Cmd{
		Name: CmdReload.CmdName,
		Help: CmdStart.Help,
		Func: func(c *ishell.Context) {
			os.Args = c.Args
		},
	})

	// TODO: debug命令
	that.Controller.AddCmd(&ishell.Cmd{
		Name: CmdDebug.CmdName,
		Help: CmdDebug.Help,
		Func: func(c *ishell.Context) {
			os.Args = c.Args
		},
	})

	// TODO: log命令
	that.Controller.AddCmd(&ishell.Cmd{
		Name: CmdLog.CmdName,
		Help: CmdLog.Help,
		Func: func(c *ishell.Context) {
			os.Args = c.Args
		},
	})
}

/*
  启动交互式Shell
*/
func (that *Keeper) StartCtrlAsClient() {
	// 主进程中启动交互式shell
	if that.Controller != nil && that.Controller.KcMode == kctrl.CtrlUnknown {
		that.Controller.InitClient()      // 初始化为Shell客户端
		that.AddCtrlCommand()             // 添加Shell命令
		that.Controller.CtrlClientStart() // 启动交互式Shell
	}
}

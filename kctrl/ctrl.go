package kctrl

import (
	"fmt"
	"net/http"

	"github.com/abiosoft/ishell/v2"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/os/gfile"
)

type CtrlMode int

const (
	CtrlUnknown CtrlMode = 0
	CtrlSever   CtrlMode = 1
	CtrlClient  CtrlMode = 2
)

type ICtrlKeeper interface {
	CtrlKeeperName() string
}

/*
  KCtrl Keeper的交互式控制器，提供交互式命令行，对本地已经启动的Keeper程序进行操作。
  可以根据需要初始化为服务端或者客户端。
*/
type KCtrl struct {
	*http.Client              // UnixSockHttp 客户端
	*gin.Engine               // UnixSockHttp 服务端
	*ishell.Shell             // 交互式Shell
	UnixSockName  string      // Unix Domain Socket名称
	UnixSockPath  string      // Unix套接字路径，默认为/tmp/[keeperName].sock
	Keeper        ICtrlKeeper // Keeper 对象
	KcMode        CtrlMode    // 控制器模式：客户端或者服务器
}

func NewKeeperCtrl(k ICtrlKeeper) *KCtrl {
	return &KCtrl{
		Keeper: k,
		KcMode: CtrlUnknown,
	}
}

// 初始化Unix套接字，服务端和客户端都有用
func (that *KCtrl) InitSockPath() {
	that.UnixSockName = fmt.Sprintf("%s.sock", that.Keeper.CtrlKeeperName())
	that.UnixSockPath = gfile.TempDir(that.UnixSockName)
}

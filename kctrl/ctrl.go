package kctrl

import (
	"errors"
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
	CtrlCurrentExecutor() string
	CtrlIsMaster() bool
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

/*
  InitSockPath 初始化Unix套接字，服务端和客户端都有用;
  主进程中使用keeperName命名；
  子进程中使用keeperName_ExecutorName命名。
*/
func (that *KCtrl) InitSockPath() {
	if that.Keeper.CtrlIsMaster() {
		// 主进程中Unix套接字
		that.UnixSockName = fmt.Sprintf("%s.sock", that.Keeper.CtrlKeeperName())
	} else {
		// 子进程中Unix套接字
		executorName := that.Keeper.CtrlCurrentExecutor()
		if executorName == "" {
			panic(errors.New("子进程中没有传入ExecutorName, 正常情况下需要传入"))
		}
		that.UnixSockPath = fmt.Sprintf("%s_%s.sock", that.Keeper.CtrlKeeperName(), executorName)
	}
	that.UnixSockPath = gfile.TempDir(that.UnixSockName)
}

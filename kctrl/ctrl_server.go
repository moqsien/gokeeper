package kctrl

import (
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/os/gfile"
	logger "github.com/moqsien/processes/logger"
)

func (that *KCtrl) InitServer() {
	if that.KcMode == CtrlUnknown {
		that.InitSockPath()
		_, err := os.Stat(that.UnixSockPath)
		if !os.IsNotExist(err) {
			// 判断socket文件是否存在，若已存在则删除
			_ = gfile.Remove(that.UnixSockPath)
		}
		that.Engine = gin.New()
		that.KcMode = CtrlSever
	}
}

func (that *KCtrl) CtrlServerStart() {
	if that.KcMode == CtrlSever {
		unixAddr, err := net.ResolveUnixAddr("unix", that.UnixSockPath)
		if err != nil {
			return
		}
		listener, err := net.ListenUnix("unix", unixAddr)
		if err != nil {
			logger.Error("listening error:", err)
			return
		}
		http.Serve(listener, that)
	} else if that.KcMode == CtrlClient {
		logger.Error("Alreay initialized to client!")
		return
	} else {
		logger.Error("Not initialized!")
		return
	}
}

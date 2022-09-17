package main

import (
	"fmt"
	"net/http"

	"github.com/abiosoft/ishell/v2"
	"github.com/gin-gonic/gin"
	"github.com/moqsien/gokeeper/kctrl"
)

/*
  kctrl使用示例
*/

type FakeKeeper struct {
	Name    string
	CServer *kctrl.KCtrl
	CClient *kctrl.KCtrl
}

var Info string = "test_info"

const (
	InfoPath = "/fake/info/"
)

func (k *FakeKeeper) CtrlKeeperName() string {
	return k.Name
}

func (k *FakeKeeper) AddRoutes() {
	if k.CServer != nil {
		k.CServer.GET(InfoPath, func(c *gin.Context) {
			c.String(http.StatusOK, Info)
		})
	}
}

func (k *FakeKeeper) AddCommand() {
	if k.CClient != nil {
		k.CClient.AddCmd(&ishell.Cmd{
			Name: "info",
			Help: "keeper info",
			Func: func(c *ishell.Context) {
				params := kctrl.ParamsContainer{}
				result, err := k.CClient.CtrlGetStr(InfoPath, params)
				if result != Info || err != nil {
					c.Printf("!!!Error Occurred!!! [Result]%s [Error] %v", result, err)
					return
				}
				fmt.Println("===Testing===")
				c.Println("+++ Show [Info]: ", result)
			},
		})
	}
}

var Default = &FakeKeeper{
	Name: "default",
}

func main() {
	Default.CServer = kctrl.NewKeeperCtrl(Default)
	Default.CServer.InitServer()
	Default.AddRoutes()
	go Default.CServer.CtrlServerStart()

	Default.CClient = kctrl.NewKeeperCtrl(Default)
	Default.CClient.InitClient()
	Default.AddCommand()
	Default.CClient.CtrlClientStart()
}

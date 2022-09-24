package keeper

import (
	"fmt"
	"os"

	"github.com/gogf/gf/os/gfile"
)

// 默认Keeper
var DefaultKeeper = NewKeeper(fmt.Sprintf("Keeper_%s", gfile.Basename(os.Args[0])))

// SetName 设置Keeper名
// 建议设置独特个性化的引用名，因为管理链接，日志目录等地方会用到它。
// 如果不设置，默认为"Keeper_xxx",启动xxx为二进制的名称
func SetName(name string) {
	DefaultKeeper.KeeperName = name
}

// Setup 设置并启动DefaultKeeper
func Setup(startFunction StartFunc) {
	DefaultKeeper.SetupStartFunc(startFunction)
}

// CloseCtrl 关闭ctrl管理功能
func CloseCtrl() {
	DefaultKeeper.CanCtrl = false
}

// Shutdown 关闭服务
func Shutdown() {
	DefaultKeeper.Shutdown()
}

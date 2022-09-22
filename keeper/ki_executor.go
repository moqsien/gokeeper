package keeper

import (
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/os/gcfg"
	ktype "github.com/moqsien/gokeeper/ktype"
	process "github.com/moqsien/processes"
)

/*
  以下方法实现Executor需要的IKeeper接口
*/
func (that *Keeper) Config() *gcfg.Config {
	return that.KConfig
}

func (that *Keeper) IsMaster() bool {
	return that.KeeperIsMaster
}

func (that *Keeper) ListOfAppsToStart() *garray.StrArray {
	return that.AppsToOperate
}

func (that *Keeper) Mode() ktype.ProcMode {
	return that.ProcMode
}

func (that *Keeper) ProcManager() *process.Manager {
	return that.Manager
}

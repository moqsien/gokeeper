package keeper

import (
	"fmt"
	"os"
	"time"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gtree"
	"github.com/gogf/gf/os/gcfg"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gutil"
	kcli "github.com/moqsien/gokeeper/kcli"
	kctrl "github.com/moqsien/gokeeper/kctrl"
	ktype "github.com/moqsien/gokeeper/ktype"
	process "github.com/moqsien/processes"
	logger "github.com/moqsien/processes/logger"
	"github.com/spf13/cobra"
)

// StartFunc keeper启动时的回调方法
type StartFunc func(k *Keeper)

// StopFunc keeper关闭前的回调方法
type StopFunc func(k *Keeper) bool

// Keeper 微服务管理者
type Keeper struct {
	*cobra.Command                        // 命令行工具
	*process.ProcManager                  // 进程管理者
	KeeperName           string           // 微服务管理者keeper的名称
	KeeperIsMaster       bool             // 是否为主进程
	KConfigPath          string           // 配置文件路径
	KConfig              *gcfg.Config     // keeper的配置信息
	ExecutorList         *gtree.AVLTree   // Executor列表
	ExecutorToOperate    string           // 要在子进程中执行的Executor
	AppsToOperate        *garray.StrArray // 需要启动或停止的App的名称列表
	PidFilePath          string           // 主进程的pid文件保存路径
	ProcMode             ktype.ProcMode   // 进程模式，MultiProcs:多进程模式；SingleProc:单进程模式
	StartFunction        StartFunc        // keeper启动方法
	StartTime            *gtime.Time      // keeper启动时间
	Exiting              bool             // keeper正在关闭
	BeforeStopFunc       StopFunc         // 服务关闭之前执行该方法
	CanCtrl              bool             // 是否开启交互式shell功能，默认true
	Controller           *kctrl.KCtrl     // 交互式shell
	// InheritAddrList      []grace.InheritAddr // 多进程模式，开启平滑重启逻辑模式下需要监听的列表
	// Graceful             *graceful.Graceful
}

// Setup 启动服务，并执行传入的启动方法
func (that *Keeper) Setup(startFunction StartFunc) {
	// // 开启ctrl命令
	// if that.openCtrl {
	// 	// ctrl命令
	// 	if len(os.Args) > 1 && os.Args[1] == "ctrl" {
	// 		os.Args = append(os.Args[0:1], os.Args[2:]...)
	// 		_ = logger.SetLevelStr("ERROR")
	// 		_, err := that.getCtrlSession()
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			os.Exit(0)
	// 		}
	// 		that.initCtrlGrumble()
	// 		err = that.grumbleApp.Run()
	// 		if err != nil {
	// 			_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)
	// 			os.Exit(1)
	// 		}
	// 		return
	// 	}
	// }
	// 如果启动进程的时候未传入任何参数，则默认使用start
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "start")
	}
	// 启动用户启动方法
	that.StartFunction = startFunction
	// 解析命令行参数，并路由参数执行逻辑
	err := that.Execute()
	if err != nil {
		that.Help()
		os.Exit(0)
	}
	//监听重启信号
	// that.graceful.graceSignal()
}

// Shutdown 主动结束进程
func (that *Keeper) Shutdown(timeout ...time.Duration) {
	// TODO: graceful shutdownSingle
	// that.graceful.shutdownSingle(timeout...)
}

// IsMutilProcModeAndInMaster 判断是否是多进程模式，并且正在主进程中执行
func (that *Keeper) IsMutilProcModeAndInMaster() bool {
	return that.ProcMode == ktype.MultiProcs && that.IsMaster()
}

// 将主进程的pid写入pid文件
func (that *Keeper) PutMasterPidInFile() {
	pid := os.Getpid()
	f, e := os.OpenFile(that.PidFilePath, os.O_WRONLY|os.O_CREATE, os.FileMode(0600))
	if e != nil {
		logger.Fatalf("os.OpenFile: %v", e)
	}
	defer func() {
		_ = f.Close()
	}()
	if e := os.Truncate(that.PidFilePath, 0); e != nil {
		logger.Fatalf("os.Truncate: %v.", e)
	}
	if _, e := fmt.Fprintf(f, "%d", pid); e != nil {
		logger.Fatalf("Unable to write pid %d to file: %s.", pid, e)
	}
	logger.Printf("写入Pid:[%d]到文件[%s]", pid, that.PidFilePath)
}

// 命令行初始化
func (that *Keeper) InitCli() {
	kcli.InitCliCmd(that)
}

// NewKeeper Keeper工厂函数
func NewKeeper(name string) *Keeper {
	svr := &Keeper{
		KeeperName:     name,
		ExecutorList:   gtree.NewAVLTree(gutil.ComparatorString, true),
		AppsToOperate:  garray.NewStrArray(true),
		ProcManager:    process.NewProcManager(),
		KeeperIsMaster: genv.GetVar(ktype.EnvIsMaster, true).Bool(),
		CanCtrl:        true,
	}
	svr.InitCli()                             // 初始化命令行
	svr.Controller = kctrl.NewKeeperCtrl(svr) // 初始化交互式命令对象
	return svr
}

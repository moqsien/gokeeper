package ktype

import "time"

/*
  一些自定义类型和公用的常量
*/

// 进程模式ProcMode
type ProcMode int

const (
	MultiProcs ProcMode = 1 // 多进程模式
	SingleProc ProcMode = 2 // 单进程模式
)

//当前进程的状态
const (
	StatusActionNone         = 0 // 初始状态
	StatusActionRestarting   = 1 // 进程在重启中
	StatusActionShuttingDown = 2 // 进程正在结束中
)

const (
	EnvIsMaster             = "ENV_MULTI_MASTER"                    // 多进程模式下，主进程的标记用环境变量名
	EnvIsChild              = "GRACEFUL_IS_CHILD"                   // 当前是否是在子进程
	ParentAddrKey           = "GRACEFUL_INHERIT_LISTEN_PARENT_ADDR" // 父进程的监听列表
	AdminActionReloadEnvKey = "GF_SERVER_RELOAD"                    // gf框架的ghttp服务平滑重启key
	MinShutdownTimeout      = 15 * time.Second                      // 进程收到结束或重启信号后，存活的最大时间
	ConfigNodeNameLogger    = "logger"
)

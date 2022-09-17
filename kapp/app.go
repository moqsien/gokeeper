package kapp

import (
	"context"

	"github.com/gogf/gf/os/gcfg"
	"github.com/gogf/gf/os/gtime"
	"github.com/moqsien/processes"
)

// AppConfig 配置文件对象
type AppConfig struct {
	*gcfg.Config
}

type IExecutor interface {
	SearchApp(name string) (IApp, bool)
}

// IApp 应用接口，对应于每个微服务应用，可以是RPC或者HTTP
type IApp interface {
	AppName() string // 获取应用名称
	Execute() error  // 启动应用，其中包装有微服务应用的业务逻辑
	Exit() error     // 关闭应用，关闭微服务应用
}

type AppBase struct {
	Executor  IExecutor       // App所属的执行器
	Context   context.Context // App专属上下文
	AppConfig *AppConfig      // App相关的配置
}

type AppContainer struct {
	App       IApp
	StartTime *gtime.Time         // APP启动时间
	StopTime  *gtime.Time         // APP关闭时间
	State     processes.ProcState // APP的运行状态，用进程状态表示
}

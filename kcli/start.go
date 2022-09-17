package kcli

import (
	"github.com/gogf/gf/os/gcfg"
	"github.com/moqsien/processes/logger"
	"github.com/spf13/cobra"
)

type IKStartCmd interface {
	SetAppsToOperate(names []string)
	SetCurrentExecutor(execName string)
	Config() *gcfg.Config
	InitLogSetting(config *gcfg.Config) error
	ParseConfig(config string)
	ParseEnv(env string)
	ParseDebug(debug bool)
	ParseDaemon(daemon bool)
	ParsePidFilePath(pidPath string)
	CheckKeeperForStart()
	RunKeeper()
}

// startCmd start命令
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start apps",
	Long:  "start apps",
}

var (
	env      string
	config   string
	executor string
	pid      string
	deamon   bool
	debug    bool
	mode     int
)

func InitStartCmd(keeper ICommand) {
	startCmd.Flags().StringVarP(&env, "env", "e", "", "环境变量，表示当前启动所在的环境,有[dev,test,product]这三种，默认是product")
	startCmd.Flags().StringVarP(&config, "config", "c", "", "指定要载入的配置文件，该参数与gf.gcfg.file参数二选一，建议使用该参数")
	startCmd.Flags().StringVarP(&executor, "executor", "x", "", "设置子进程需要启动的Executor名称，默认为空")
	startCmd.Flags().StringVarP(&pid, "pid", "p", "", "设置pid文件的地址，默认是/tmp/[keeperName].pid")
	startCmd.Flags().BoolVarP(&deamon, "deamon", "d", false, "使用守护进程模式启动")
	startCmd.Flags().BoolVarP(&debug, "debug", "d", false, "是否开启debug 默认debug=true")
	startCmd.Flags().IntVarP(&mode, "mode", "m", 0, "进程模型，0表示单进程模型，1表示多进程模型")
	startCmd.Run = func(c *cobra.Command, args []string) {
		/*
		  以下为参数解析
		*/
		keeper.SetAppsToOperate(args)
		if kconfig, err := c.Flags().GetString("config"); err == nil {
			keeper.ParseConfig(kconfig)
		}
		if kexecutor, err := c.Flags().GetString("executor"); err == nil {
			keeper.SetCurrentExecutor(kexecutor) // 开启子进程时传入要运行的ExecutorName
		}
		if kenv, err := c.Flags().GetString("env"); err == nil {
			keeper.ParseEnv(kenv)
		}
		if kdebug, err := c.Flags().GetBool("debug"); err == nil {
			keeper.ParseDebug(kdebug)
		}
		if kdeamon, err := c.Flags().GetBool("deamon"); err == nil {
			keeper.ParseDaemon(kdeamon)
		}
		if kpid, err := c.Flags().GetString("pid"); err == nil {
			keeper.ParsePidFilePath(kpid)
		}

		//初始化日志配置
		if e := keeper.InitLogSetting(keeper.Config()); e != nil {
			logger.Fatalf("error:%v", e)
		}
		keeper.CheckKeeperForStart()

		keeper.RunKeeper()
	}
	keeper.AddCommand(startCmd)
}

package kutils

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/gogf/gf/os/gcfg"
	"github.com/gogf/gf/os/gfile"
	logger "github.com/moqsien/processes/logger"
)

// 以守护进程运行
func Demonize(config *gcfg.Config) error {
	// 判断是否需要后台运行
	daemon := config.GetBool("Daemon", false)
	if !daemon {
		return nil
	}

	if syscall.Getppid() == 1 {
		return nil
	}
	// 将命令行参数中执行文件路径转换成可用路径
	filePath := gfile.SelfPath()
	logger.Infof("Starting %s: ", filePath)
	arg0, e := exec.LookPath(filePath)
	if e != nil {
		return e
	}
	argv := make([]string, 0, len(os.Args))
	for _, arg := range os.Args {
		if arg == "--daemon" || arg == "-d" {
			continue
		}
		argv = append(argv, arg)
	}
	cmd := exec.Command(arg0, argv[1:]...)
	cmd.Env = os.Environ()
	// 将其他命令传入生成出的进程
	cmd.Stdin = os.Stdin // 给新进程设置文件描述符，可以重定向到文件中
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start() // 开始执行新进程，不等待新进程退出
	if err != nil {
		return err
	}
	os.Exit(0)
	return nil
}

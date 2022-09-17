package kcli

import (
	"fmt"

	"github.com/spf13/cobra"
)

/*
  Keeper的命令行参数解析
*/

type ICommand interface {
	SetRootCommand(c *cobra.Command)
	AddCommand(cmds ...*cobra.Command)
	Execute() error
	IKStartCmd
	IkStopCmd
	IkVersionCmd
}

var RootCmd = &cobra.Command{
	Use:   "keeper [keeper]",
	Short: "microservice keeper",
	Long:  "use -h or --help to see how to use",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This is Keeper")
	},
}

func InitCliCmd(kPtr ICommand) {
	kPtr.SetRootCommand(RootCmd) // kPtr为Keeper指针
	InitStartCmd(kPtr)
	InitStopCmd(kPtr)
	InitReloadCmd(kPtr)
	InitRQuitCmd(kPtr)
	InitVersionCmd(kPtr)
}

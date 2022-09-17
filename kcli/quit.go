package kcli

import "github.com/spf13/cobra"

var quitCmd = &cobra.Command{
	Use:   "quit",
	Short: "quit apps",
	Long:  "quit apps",
}

func InitRQuitCmd(keeper ICommand) {
	quitCmd.Run = func(c *cobra.Command, args []string) {
		keeper.SetAppsToOperate(args)
		if pid, err := c.Flags().GetString("pid"); err == nil {
			keeper.ParsePidFilePath(pid)
		}
		keeper.StopKeeper("quit")
	}
	stopCmd.Flags().StringVarP(&pid, "pid", "p", "", "设置pid文件的地址，默认是/tmp/[keeperName].pid")
	keeper.AddCommand(quitCmd)
}

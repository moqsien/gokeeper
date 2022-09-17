package kcli

import "github.com/spf13/cobra"

type IkStopCmd interface {
	StopKeeper(sig string)
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop apps",
	Long:  "stop apps",
}

func InitStopCmd(keeper ICommand) {
	stopCmd.Run = func(c *cobra.Command, args []string) {
		keeper.SetAppsToOperate(args)
		if pid, err := c.Flags().GetString("pid"); err == nil {
			keeper.ParsePidFilePath(pid)
		}
		keeper.StopKeeper("stop")
	}
	stopCmd.Flags().StringVarP(&pid, "pid", "p", "", "设置pid文件的地址，默认是/tmp/[keeperName].pid")
	keeper.AddCommand(stopCmd)
}

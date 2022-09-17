package kcli

import "github.com/spf13/cobra"

var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "reload apps",
	Long:  "reload apps",
}

func InitReloadCmd(keeper ICommand) {
	reloadCmd.Run = func(c *cobra.Command, args []string) {
		keeper.SetAppsToOperate(args)
		if pid, err := c.Flags().GetString("pid"); err == nil {
			keeper.ParsePidFilePath(pid)
		}
		keeper.StopKeeper("reload")
	}
	stopCmd.Flags().StringVarP(&pid, "pid", "p", "", "设置pid文件的地址，默认是/tmp/[keeperName].pid")
	keeper.AddCommand(reloadCmd)
}

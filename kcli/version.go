package kcli

import (
	"os"

	"github.com/spf13/cobra"
)

type IkVersionCmd interface {
	Version()
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show version info",
	Long:  "show version info",
}

func InitVersionCmd(keeper ICommand) {
	versionCmd.Run = func(c *cobra.Command, args []string) {
		keeper.Version()
		os.Exit(0)
	}
	keeper.AddCommand(versionCmd)
}

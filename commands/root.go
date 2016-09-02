package commands

import (
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "passage",
	Short: "Passage - SSH tunnels on steroids",
}

func init() {
	RootCmd.AddCommand(NewServerCommand().Command())
	RootCmd.AddCommand(NewGetCommand().Command())
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "passage",
	Short: "Passage - SSH tunnels on steroids",
}

func init() {
	RootCmd.AddCommand(NewServerCommand().Command())
	RootCmd.AddCommand(NewListenCommand().Command())
	RootCmd.AddCommand(NewGetCommand().Command())
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

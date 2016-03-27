package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "passage",
	Short: "A brief description of your application",
	Long:  `A longer description`,
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

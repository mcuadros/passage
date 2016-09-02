package main

import (
	"os"

	"github.com/mcuadros/passage/commands"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

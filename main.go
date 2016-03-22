package main

import (
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/mcuadros/passage/commands"
)

func main() {
	parser := flags.NewNamedParser("passage", flags.Default)
	parser.AddCommand("listen", "", "", &commands.ListenCommand{})

	_, err := parser.Parse()
	if err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrCommandRequired {
			parser.WriteHelp(os.Stdout)
		}

		os.Exit(1)
	}
}

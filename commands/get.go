package commands

import (
	"fmt"
	"net/rpc"

	"github.com/spf13/cobra"
)

type GetCommand struct {
	RPCAddr string
}

func NewGetCommand() *GetCommand {
	return &GetCommand{}
}

func (c *GetCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [passage-name]",
		Short: "A brief description of your command",
		Long:  `A longeraaaa description .`,
		RunE:  c.Execute,
	}

	cmd.Flags().StringVar(&c.RPCAddr, "rpc-addr", "/tmp/passage.sock", "passage rpc server address, is an unix socket.")

	return cmd
}

func (l *GetCommand) Execute(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("invalid args: %q", args)
	}

	rpcClient, err := rpc.Dial("unix", l.RPCAddr)
	if err != nil {
		return err
	}

	var reply string
	err = rpcClient.Call("Server.Addr", args[0], &reply)
	if err != nil {
		return err
	}

	fmt.Printf(reply)
	return nil
}

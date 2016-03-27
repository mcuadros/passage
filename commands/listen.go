package commands

import (
	"fmt"
	"os/user"

	"github.com/mcuadros/passage/core"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

type ListenCommand struct {
	User   string
	Addr   Addr
	Server ServerAddr
	Remote Remote
}

func NewListenCommand() *ListenCommand {
	return &ListenCommand{}
}

func (c *ListenCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listen [ssh-server] [remote-addr]",
		Short: "A brief description of your command",
		Long:  `A longeraaaa description .`,
		RunE:  c.Execute,
	}

	cmd.Flags().StringVar(&c.User, "user", "", "user used in the ssh connection, if empty the current one is used")
	cmd.Flags().Var(&c.Addr, "addr", "local bind address")

	return cmd
}

func (l *ListenCommand) Execute(cmd *cobra.Command, args []string) error {
	if err := l.loadArgs(args); err != nil {
		return err
	}

	if err := l.setUser(); err != nil {
		return err
	}

	c := core.NewSSHConnection(
		l.Server,
		&ssh.ClientConfig{
			User: l.User,
			Auth: []ssh.AuthMethod{
				core.SSHAgent(),
			},
		},
	)

	p := core.NewPassage(c, l.Remote)
	if err := p.Start(l.Addr); err != nil {
		return err
	}

	fmt.Printf("(%s@%s)-[%s]->(%s)\n", l.User, l.Server, l.Remote, p)

	select {}

	return nil
}

func (l *ListenCommand) loadArgs(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing arguments: [ssh-server] [remote-addr]")
	}

	if len(args) != 2 {
		return fmt.Errorf("invalid arguments: %s", args)
	}

	if err := l.Server.Set(args[0]); err != nil {
		return err
	}

	if err := l.Remote.Set(args[1]); err != nil {
		return err
	}

	fmt.Println(l.Server, l.Remote)

	return nil
}

func (l *ListenCommand) setUser() error {
	if l.User != "" {
		return nil
	}

	u, err := user.Current()
	if err != nil {
		return err
	}

	l.User = u.Username
	return nil
}

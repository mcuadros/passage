package commands

import (
	"fmt"
	"os"
	"os/user"

	"github.com/mcuadros/passage/core"
	"golang.org/x/crypto/ssh"
)

type ListenCommand struct {
	User string `long:"user" description:"ssh server user" default:""`
	Addr Addr   `long:"addr" description:"local bind address" default:"localhost:0"`
	Args struct {
		Server ServerAddr `positional-arg-name:"server" description:"." required:"true"`
		Remote Remote     `positional-arg-name:"remote" description:"." required:"true"`
	} `positional-args:"yes"`
}

func (l *ListenCommand) Execute(args []string) error {
	if err := l.setUser(); err != nil {
		return err
	}

	c := core.NewSSHConnection(
		l.Args.Server,
		&ssh.ClientConfig{
			User: l.User,
			Auth: []ssh.AuthMethod{
				SSHAgent(),
			},
		},
	)

	p := core.NewPassage(c, l.Args.Remote)
	err := p.Start(l.Addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("(%s@%s)-[%s]->(%s)\n", l.User, l.Args.Server, l.Args.Remote, p)

	select {}

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

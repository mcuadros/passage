package commands

import (
	"fmt"
	"net"
	"strings"

	"github.com/mcuadros/passage/core"
)

const rpcAddrDefault = "/tmp/passage.sock"

type ServerAddr struct {
	Addr
}

func (a *ServerAddr) Type() string { return "server-addr" }
func (a *ServerAddr) Set(value string) error {
	if strings.Index(value, ":") == -1 {
		value = fmt.Sprintf("%s:22", value)
	}

	return a.Addr.Set(value)
}

type Addr struct {
	core.Addr
}

func (a *Addr) Type() string { return "addr" }
func (a *Addr) Set(value string) error {
	addr, err := net.ResolveTCPAddr("tcp", value)
	if err != nil {
		return err
	}

	a.Addr = core.Addr{addr}
	return nil
}

func (a Addr) String() string {
	if a.Addr.Addr == nil {
		return "localhost:0"
	}

	return a.Addr.String()
}

type Remote struct {
	core.Remote
}

func (r *Remote) Type() string { return "remote" }

func (r *Remote) Set(value string) error {
	slash := strings.Split(value, "/")

	network := "tcp"
	if len(slash) == 2 {
		network = slash[1]
	}

	if slash[0][0] == ':' {
		slash[0] = slash[0][1:]
	}

	if strings.Count(slash[0], ":") == 0 {
		r.Remote = core.NewRemote(network, fmt.Sprintf("127.0.0.1:%s", slash[0]))
		return nil
	}

	equal := strings.Split(slash[0], "=")
	if len(equal) == 1 {
		r.Remote = core.NewRemote(network, slash[0])
		return nil
	}

	dots := strings.Split(equal[1], ":")

	switch equal[0] {
	case "container":
		r.Remote = core.NewContainerRemote(network, dots[0], dots[1])
		return nil
	}

	if r.Remote == nil {
		return fmt.Errorf("invalid remote format: %s", value)

	}

	return nil
}

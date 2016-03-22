package commands

import (
	"fmt"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	"github.com/mcuadros/passage/core"
)

type ServerAddr struct {
	Addr
}

func (a *ServerAddr) UnmarshalFlag(value string) error {
	if strings.Index(value, ":") == -1 {
		value = fmt.Sprintf("%s:22", value)
	}

	return a.Addr.UnmarshalFlag(value)
}

type Addr struct {
	core.Addr
}

func (a *Addr) UnmarshalFlag(value string) error {
	addr, err := net.ResolveTCPAddr("tcp", value)
	if err != nil {
		return err
	}

	a.Addr = core.Addr{addr}
	return nil
}

type Remote struct {
	core.Remote
}

//<kind>:address:port/proto
func (r *Remote) UnmarshalFlag(value string) error {
	network, err := r.getNetwork(value)
	if err != nil {
		return err
	}

	if value[0] == ':' {
		value = value[1:]
	}

	dots := strings.Split(value, ":")
	slash := strings.Split(dots[len(dots)-1], "/")

	switch len(dots) {
	case 1:
		r.Remote = core.NewLocalhostRemote(network, slash[0])
	case 2:
		if err := r.buildRemote(network, dots[0], slash[0]); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid remote format: %s", value)
	}

	return nil
}

func (r *Remote) buildRemote(network, address, port string) error {
	parts := strings.Split(address, "=")
	switch len(parts) {
	case 1:
		r.Remote = core.NewRemote(network, address, port)
	case 2:
		switch parts[0] {
		case "container":
			r.Remote = core.NewContainerRemote(network, parts[1], port)
		default:
			return fmt.Errorf("invalid remote kind: %s", parts[0])
		}
	default:
		return fmt.Errorf("invalid remote address format: %s", address)
	}

	return nil
}

func (r *Remote) getNetwork(value string) (string, error) {
	parts := strings.Split(value, "/")
	switch len(parts) {
	case 1:
		return "tcp", nil
	case 2:
		return parts[1], nil
	default:
		return "", fmt.Errorf("invalid remote format: %s", value)
	}
}

func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

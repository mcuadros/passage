package core

import (
	"fmt"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type Addr struct {
	net.Addr
}

func (a *Addr) Address() string {
	s := strings.Split(a.Addr.String(), ":")
	return s[0]
}

func (a *Addr) Port() string {
	s := strings.Split(a.Addr.String(), ":")
	return s[1]
}

var supportedNetworks = map[string]bool{
	"udp": true, "udp4": true, "udp6": true,
	"tcp": true, "tcp4": true, "tcp6": true,
}

func MustResolveAddr(network, address string) *Addr {
	if _, ok := supportedNetworks[network]; !ok {
		panic(fmt.Sprintf("invalid network: %s", network))
	}

	var a net.Addr
	var err error
	switch network[:3] {
	case "tcp":
		a, err = net.ResolveTCPAddr(network, address)
	case "udp":
		a, err = net.ResolveUDPAddr(network, address)
	}

	if err != nil {
		panic(err)
	}

	return &Addr{a}
}

func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}

	return nil
}

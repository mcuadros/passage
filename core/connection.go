package core

import (
	"fmt"
	"io"
	"net"
	"sync"

	"golang.org/x/crypto/ssh"
)

type SSHConnection interface {
	Tunnel(c net.Conn, a net.Addr) error
	Conn(a net.Addr) (net.Conn, error)
	fmt.Stringer
}

type sshConnection struct {
	a net.Addr
	c *ssh.ClientConfig

	client *ssh.Client
}

func NewSSHConnection(a net.Addr, c *ssh.ClientConfig) SSHConnection {
	return &sshConnection{a: a, c: c}
}

func (s *sshConnection) Tunnel(c net.Conn, a net.Addr) error {
	r, err := s.Conn(a)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	copyConn := func(writer, reader net.Conn) {
		defer wg.Done()
		if _, err := io.Copy(writer, reader); err != nil {
			fmt.Printf("io.Copy error: %s", err)
		}
	}

	wg.Add(2)
	go copyConn(c, r)
	go copyConn(r, c)
	wg.Wait()

	return nil
}

func (s *sshConnection) Conn(a net.Addr) (net.Conn, error) {
	if err := s.dialServerConnection(); err != nil {
		return nil, err
	}

	return s.dialRemoteConnection(a)
}

func (c *sshConnection) dialServerConnection() error {
	if c.client != nil {
		return nil
	}

	var err error
	c.client, err = ssh.Dial(c.a.Network(), c.a.String(), c.c)
	if err != nil {
		return err
	}

	return nil
}

func (c *sshConnection) dialRemoteConnection(a net.Addr) (net.Conn, error) {
	return c.client.Dial(a.Network(), a.String())
}

func (c *sshConnection) String() string {
	return fmt.Sprintf("%s@%s", c.c.User, c.a)
}

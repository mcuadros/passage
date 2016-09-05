package core

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHConnection interface {
	Tunnel(c net.Conn, a net.Addr) error
	Conn(a net.Addr) (net.Conn, error)
	Config() *ssh.ClientConfig
	fmt.Stringer
}

type sshConnection struct {
	a          net.Addr
	c          *ssh.ClientConfig
	maxRetries int

	connected bool
	client    *ssh.Client
}

func NewSSHConnection(a net.Addr, c *ssh.ClientConfig, retries int) SSHConnection {
	return &sshConnection{a: a, c: c, maxRetries: retries}
}

func (s *sshConnection) Config() *ssh.ClientConfig {
	return s.c
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

func (c *sshConnection) Conn(a net.Addr) (net.Conn, error) {
	conn, err := c.dialRemoteConnection(a)
	if err == nil {
		return conn, nil
	}

	c.connected = false
	var retries int
	for range time.Tick(5 * time.Second) {
		conn, err := c.dialRemoteConnection(a)
		if err == nil {
			return conn, nil
		}

		retries++
		if retries > c.maxRetries {
			return nil, fmt.Errorf("%s, after %d retries", err, retries-1)
		}
	}

	panic("unrechable")
}

func (c *sshConnection) dialRemoteConnection(a net.Addr) (net.Conn, error) {
	if err := c.dialServerConnection(); err != nil {
		return nil, err
	}

	conn, err := c.client.Dial(a.Network(), a.String())
	if err != nil {
		return nil, fmt.Errorf("error dialing remote: %s", err)
	}

	return conn, nil
}

func (c *sshConnection) dialServerConnection() error {
	if c.connected {
		return nil
	}

	c.c.Timeout = time.Second * 5
	var err error
	c.client, err = ssh.Dial(c.a.Network(), c.a.String(), c.c)
	if err != nil {
		return fmt.Errorf("error dialing server: %s", err)
	}

	c.connected = true
	return nil
}

func (c *sshConnection) String() string {
	return fmt.Sprintf("%s@%s", c.c.User, c.a)
}

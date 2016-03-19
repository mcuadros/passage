package passage

import (
	"fmt"
	"io"
	"net"
	"sync"

	"golang.org/x/crypto/ssh"
)

type SSHServer struct {
	a net.Addr
	c *ssh.ClientConfig

	client *ssh.Client
}

func NewSSHServer(a net.Addr, c *ssh.ClientConfig) *Server {
	return &SSHServer{a: a, c: c}
}

func (s *SSHServer) Tunnel(c net.Conn, a net.Addr) error {
	if err := s.dialServerConnection(); err != nil {
		return err
	}

	r, err := s.dialRemoteConnection(a)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	copyConn := func(writer, reader net.Conn) {
		defer wg.Done()
		_, err := io.Copy(writer, reader)
		if err != nil {
			fmt.Printf("io.Copy error: %s", err)
		}
	}

	wg.Add(2)
	go copyConn(c, r)
	go copyConn(r, c)
	wg.Wait()

	return nil
}

func (s *SSHServer) dialServerConnection() error {
	if s.client != nil {
		return nil
	}

	var err error
	s.client, err = ssh.Dial(s.a.Network(), s.a.String(), s.c)
	if err != nil {
		return err
	}

	return nil
}

func (s *SSHServer) dialRemoteConnection(a net.Addr) (net.Conn, error) {
	return s.client.Dial(a.Network(), a.String())
}

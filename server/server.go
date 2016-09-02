package server

import (
	"fmt"
	"net"
	"strings"

	"github.com/mcuadros/passage/core"

	"golang.org/x/crypto/ssh"
	"gopkg.in/inconshreveable/log15.v2"
)

type Server struct {
	c *Config

	servers  map[string]core.SSHConnection
	passages map[string]*core.Passage
}

func NewServer() *Server {
	return &Server{
		servers:  make(map[string]core.SSHConnection, 0),
		passages: make(map[string]*core.Passage, 0),
	}
}

func (s *Server) Load(c *Config) error {
	if err := c.Validate(); err != nil {
		return err
	}

	var loadedServers, loadedPassages []string

	for name, server := range c.Servers {
		loadedPassage, err := s.loadSSHConnection(name, server)
		if err != nil {
			return err
		}

		loadedServers = append(loadedServers, name)
		loadedPassages = append(loadedPassages, loadedPassage...)
	}

	s.cleanServers(loadedServers)
	s.cleanPassages(loadedPassages)
	return nil
}

func (s *Server) loadSSHConnection(name string, config *SSHServerConfig) ([]string, error) {
	c, err := s.buildSSHConnection(config)
	if err != nil {
		return nil, err
	}

	if _, ok := s.servers[name]; !ok {
		s.servers[name] = c
	}

	loadedPassages, err := s.loadPassages(s.servers[name], config)
	if err != nil {
		return loadedPassages, err
	}

	return loadedPassages, nil
}

func (s *Server) buildSSHConnection(config *SSHServerConfig) (core.SSHConnection, error) {
	a, err := net.ResolveTCPAddr("tcp", config.Address)
	if err != nil {
		return nil, err
	}

	return core.NewSSHConnection(a, &ssh.ClientConfig{
		User:    config.User,
		Timeout: config.Timeout,
		Auth: []ssh.AuthMethod{
			core.SSHAgent(),
		},
	}, config.Retries), nil
}

func (s *Server) loadPassages(c core.SSHConnection, config *SSHServerConfig) ([]string, error) {
	var loadedPassages []string
	for name, p := range config.Passages {
		if err := s.loadPassage(c, name, p); err != nil {
			return loadedPassages, err
		}

		loadedPassages = append(loadedPassages, name)
	}

	return loadedPassages, nil

}

func (s *Server) loadPassage(c core.SSHConnection, name string, config *PassageConfig) error {
	r, err := s.buildRemote(config)
	if err != nil {
		return err
	}

	a, err := net.ResolveTCPAddr("tcp", config.Local)
	if err != nil {
		return err
	}

	if _, ok := s.passages[name]; ok {
		return nil
	}

	s.passages[name] = core.NewPassage(c, r)
	if err := s.passages[name].Start(a); err != nil {
		return err
	}

	log15.Info(
		"new passage created",
		"name", name, "ssh", c, "remote", r, "addr", s.passages[name].Addr(),
	)

	return nil
}

func (s *Server) buildRemote(config *PassageConfig) (core.Remote, error) {
	switch config.Type {
	case "tcp":
		return core.NewRemote("tcp", config.Address), nil
	case "container":
		return core.NewContainerRemote("tcp", config.Container, config.Port), nil
	}

	return nil, fmt.Errorf("invalid remote type: %q", config.Type)
}

func (s *Server) cleanServers(loadedServers []string) {
	for k := range s.servers {
		if !contains(loadedServers, k) {
			delete(s.servers, k)
		}
	}
}

func (s *Server) cleanPassages(loadedPassages []string) {
	var removed []string
	for k := range s.passages {
		if !contains(loadedPassages, k) {
			delete(s.passages, k)
			removed = append(removed, k)
		}
	}

	if len(removed) == 0 {
		return
	}

	log15.Debug("removed passages", "names", removed)
}

func (s *Server) Close() error {
	for _, p := range s.passages {
		if err := p.Close(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) String() string {
	var out []string
	for _, p := range s.passages {
		out = append(out, p.String())
	}

	return strings.Join(out, "\n")
}

func contains(haystack []string, needle string) bool {
	for _, e := range haystack {
		if e == needle {
			return true
		}
	}

	return false
}

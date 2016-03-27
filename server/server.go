package server

import (
	"fmt"
	"net"

	"gopkg.in/inconshreveable/log15.v2"

	"golang.org/x/crypto/ssh"

	"github.com/mcuadros/go-defaults"
	"github.com/mcuadros/passage/core"
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
	var loadedServers, loadedPassages []string

	for _, server := range c.Servers {
		loadedServer, loadedPassage, err := s.loadSSHConnection(&server)
		if err != nil {
			return err
		}

		loadedServers = append(loadedServers, loadedServer)
		loadedPassages = append(loadedPassages, loadedPassage...)
	}

	s.cleanServers(loadedServers)
	s.cleanPassages(loadedPassages)
	return nil
}

func (s *Server) loadSSHConnection(config *SSHServerConfig) (string, []string, error) {
	defaults.SetDefaults(config)

	c, err := s.buildSSHConnection(config)
	if err != nil {
		return "", nil, err
	}

	key := c.String()
	if _, ok := s.servers[key]; !ok {
		s.servers[key] = c
	}

	loadedPassages, err := s.loadPassages(s.servers[key], config)
	if err != nil {
		return key, loadedPassages, err
	}

	return key, loadedPassages, nil
}

func (s *Server) buildSSHConnection(config *SSHServerConfig) (core.SSHConnection, error) {
	a, err := net.ResolveTCPAddr("tcp", config.Address)
	if err != nil {
		return nil, err
	}

	return core.NewSSHConnection(a, &ssh.ClientConfig{
		User: config.User,
		Auth: []ssh.AuthMethod{
			core.SSHAgent(),
		},
	}), nil
}

func (s *Server) loadPassages(c core.SSHConnection, config *SSHServerConfig) ([]string, error) {
	var loadedPassages []string
	for _, p := range config.Passages {
		loaded, err := s.loadPassage(c, &p)
		if err != nil {
			return loadedPassages, err
		}

		loadedPassages = append(loadedPassages, loaded)
	}

	return loadedPassages, nil

}

func (s *Server) loadPassage(c core.SSHConnection, config *PassageConfig) (string, error) {
	defaults.SetDefaults(config)

	r, err := s.buildRemote(&config.Remote)
	if err != nil {
		return "", err
	}

	a, err := net.ResolveTCPAddr("tcp", config.Local)
	if err != nil {
		return "", err
	}

	p := core.NewPassage(c, r)

	key := config.Name
	if key == "" {
		key = p.String()
	}

	if _, ok := s.passages[key]; ok {
		return key, nil
	}

	if err := p.Start(a); err != nil {
		return key, err
	}

	s.passages[key] = p

	logArgs := log15.Ctx{"ssh": c, "remote": r, "addr": p.Addr()}
	if config.Name != "" {
		logArgs["name"] = config.Name
	}

	log15.Info("new passage created", logArgs)

	return key, nil
}

func (s *Server) buildRemote(config *RemoteConfig) (core.Remote, error) {
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
	for _, p := range s.passages {
		fmt.Println(p)
	}

	return "foo"
}

func contains(haystack []string, needle string) bool {
	for _, e := range haystack {
		if e == needle {
			return true
		}
	}

	return false
}

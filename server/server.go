package server

import (
	"crypto/sha1"
	"fmt"
	"net"
	"strings"

	"github.com/mcuadros/passage/core"

	"golang.org/x/crypto/ssh"
	"gopkg.in/inconshreveable/log15.v2"
)

type Server struct {
	c *Config
	f fingerprints

	servers  map[string]core.SSHConnection
	passages map[string]*core.Passage
}

func NewServer() *Server {
	return &Server{
		f:        make(fingerprints),
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

	if s.f.IsNewSSHServer(name, config) {
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

	agent, err := core.SSHAgent()
	if err != nil {
		return nil, err
	}

	return core.NewSSHConnection(a, &ssh.ClientConfig{
		User:    config.User,
		Timeout: config.Timeout,
		Auth:    []ssh.AuthMethod{agent},
	}, config.Retries), nil
}

func (s *Server) loadPassages(c core.SSHConnection, config *SSHServerConfig) ([]string, error) {
	var loadedPassages []string
	for name, p := range config.Passages {
		if err := s.loadPassage(c, config, name, p); err != nil {
			return loadedPassages, err
		}

		loadedPassages = append(loadedPassages, name)
	}

	return loadedPassages, nil

}

func (s *Server) loadPassage(
	c core.SSHConnection, sc *SSHServerConfig, name string, config *PassageConfig,
) error {
	r, err := s.buildRemote(config)
	if err != nil {
		return err
	}

	a, err := net.ResolveTCPAddr("tcp", config.Local)
	if err != nil {
		return err
	}

	if !s.f.IsNewPassage(name, sc, config) {
		return nil
	}

	if _, ok := s.passages[name]; ok {
		if err := s.passages[name].Close(); err != nil {
			return err
		}
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

type fingerprints map[string][20]byte

func (fp *fingerprints) IsNewSSHServer(id string, c *SSHServerConfig) bool {
	hash := fp.fpSSHServer(c)
	if (*fp)[id] == hash {
		return false
	}

	(*fp)[id] = hash
	return true
}

func (fp *fingerprints) fpSSHServer(c *SSHServerConfig) [20]byte {
	payload := fmt.Sprintf("%s,%d,%s,%s", c.Address, c.Retries, c.Timeout, c.User)
	return sha1.Sum([]byte(payload))
}

func (fp *fingerprints) IsNewPassage(id string, s *SSHServerConfig, p *PassageConfig) bool {
	hash := fp.fpPassage(s, p)
	if (*fp)[id] == hash {
		return false
	}

	(*fp)[id] = hash
	return true
}

func (fp *fingerprints) fpPassage(s *SSHServerConfig, p *PassageConfig) [20]byte {
	payload := fmt.Sprintf("%v,%s", p, fp.fpSSHServer(s))
	return sha1.Sum([]byte(payload))
}

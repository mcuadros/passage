package passage

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

var supportedNetworks = map[string]bool{
	"udp": true, "udp4": true, "udp6": true,
	"tcp": true, "tcp4": true, "tcp6": true,
}

type Remote interface {
	Addr(SSHConnection) (net.Addr, error)
}

type addressRemote struct {
	network string
	address string
	port    string
}

func NewRemote(network, address, port string) Remote {
	if _, ok := supportedNetworks[network]; !ok {
		return nil
	}

	return &addressRemote{
		network: network,
		address: address,
		port:    port,
	}
}

func NewLocalhostRemote(network, port string) Remote {
	return NewRemote(network, "localhost", port)
}

func (r *addressRemote) Addr(SSHConnection) (net.Addr, error) {
	return net.ResolveTCPAddr(
		r.network,
		net.JoinHostPort(r.address, r.port),
	)
}

type containerRemote struct {
	container string
	port      string
}

func NewContainerRemote(container, port string) Remote {
	return &containerRemote{
		container: container,
		port:      port,
	}
}

func (r *containerRemote) Addr(s SSHConnection) (net.Addr, error) {
	c := r.buildClient(s)
	ip, err := r.getContainerIP(c)
	if err != nil {
		return nil, err
	}

	return net.ResolveTCPAddr("tcp", net.JoinHostPort(ip, r.port))
}

func (r *containerRemote) buildClient(s SSHConnection) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			ResponseHeaderTimeout: time.Second * 2,
			Dial: func(network, address string) (net.Conn, error) {
				a, err := net.ResolveTCPAddr(network, address)
				if err != nil {
					return nil, err
				}

				return s.Conn(a)
			},
		},
	}
}

func (r *containerRemote) getContainerIP(c *http.Client) (string, error) {
	l, err := r.getContainers(c)
	if err != nil {
		return "", err
	}

	container, err := r.matchContainer(l)
	if err != nil {
		return "", err
	}

	n, ok := container.NetworkSettings.Networks["bridge"]
	if !ok {
		return "", fmt.Errorf("container: not supported networks")
	}

	return n.IPAddress, nil
}

type container struct {
	Names           []string
	NetworkSettings struct {
		Networks map[string]struct {
			IPAddress string
		}
	}
}

func (r *containerRemote) matchContainer(l []*container) (*container, error) {
	var container *container
	for _, c := range l {
		for _, name := range c.Names {
			if name == fmt.Sprintf("/%s", r.container) {
				container = c
			}
		}
	}

	if container == nil {
		return nil, fmt.Errorf("container %q not found", r.container)
	}

	return container, nil
}

func (r *containerRemote) getContainers(c *http.Client) ([]*container, error) {
	req, err := http.NewRequest("GET", "http://localhost:2375/containers/json", nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(res.Body)
	defer res.Body.Close()

	result := []*container{}
	err = dec.Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

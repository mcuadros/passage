package server

import "gopkg.in/yaml.v1"

type Config struct {
	Servers []SSHServerConfig
}

type SSHServerConfig struct {
	User     string `default:"root"`
	Address  string
	Passages []PassageConfig
}

type PassageConfig struct {
	Remote RemoteConfig
	Local  string
}

type RemoteConfig struct {
	Type      string `default:"tcp"`
	Address   string
	Container string
	Port      string
}

func (c *Config) Marshal() ([]byte, error) {
	return yaml.Marshal(c)
}

func (c *Config) Unmarshal(in []byte) error {
	return yaml.Unmarshal(in, c)
}

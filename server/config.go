package server

import (
	"fmt"
	"strings"

	"github.com/mcuadros/go-defaults"
	"gopkg.in/yaml.v1"
)

type Config struct {
	Servers map[string]*SSHServerConfig
}

func (c *Config) Validate() error {
	defaults.SetDefaults(c)

	if len(c.Servers) == 0 {
		return fmt.Errorf("invalid empty config")
	}

	var errs []error
	for name, sc := range c.Servers {
		if err := sc.validate(name); len(err) != 0 {
			errs = append(errs, err...)
		}
	}

	if len(errs) != 0 {
		return &ConfigError{errs}
	}

	return nil
}

type SSHServerConfig struct {
	User     string `default:"root"`
	Address  string
	Passages map[string]*PassageConfig
}

func (c *SSHServerConfig) validate(name string) []error {
	defaults.SetDefaults(c)
	var errs []error

	if c.User == "" {
		errs = append(errs, fmt.Errorf("ssh server %q: user cannot be empty", name))
	}

	if c.Address == "" {
		errs = append(errs, fmt.Errorf("ssh server %q: address cannot be empty", name))
	}

	if len(c.Passages) == 0 {
		errs = append(errs, fmt.Errorf("ssh server %q: passages cannot be empty", name))
	}

	for name, pc := range c.Passages {
		if err := pc.validate(name); len(err) != 0 {
			errs = append(errs, err...)
		}
	}

	return errs
}

type PassageConfig struct {
	Type      string `default:"tcp"`
	Address   string
	Container string
	Port      string
	Local     string `default:"127.0.0.1:0"`
}

func (c *PassageConfig) validate(name string) []error {
	defaults.SetDefaults(c)

	var errs []error
	if c.Local == "" {
		errs = append(errs, fmt.Errorf("passage %q: local cannot be empty", name))
	}

	if valid := PassageConfigValidTypes[c.Type]; !valid {
		errs = append(errs, fmt.Errorf("passage %q: invalid remote type %q", name, c.Type))
	}

	return errs
}

var PassageConfigValidTypes = map[string]bool{"tcp": true, "container": true}

func (c *Config) Marshal() ([]byte, error) {
	return yaml.Marshal(c)
}

func (c *Config) Unmarshal(in []byte) error {
	return yaml.Unmarshal(in, c)
}

type ConfigError struct {
	Errors []error
}

func (err *ConfigError) Error() string {
	var s []string
	for _, e := range err.Errors {
		s = append(s, fmt.Sprintf("\t%s", e))
	}

	return fmt.Sprintf(
		"Invalid configuration, found %d error(s):\n%s",
		len(s), strings.Join(s, "\n"),
	)
}

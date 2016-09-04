package server

import (
	"fmt"
	"os/user"
	"strings"
	"time"

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

	errs = append(errs, c.validatePassageNames()...)
	if len(errs) != 0 {
		return &ConfigError{errs}
	}

	return nil
}

func (c *Config) validatePassageNames() []error {
	seen := map[string]bool{}
	var errs []error

	for server, s := range c.Servers {
		for n := range s.Passages {
			if seen[n] {
				errs = append(errs,
					fmt.Errorf("ssh server %q: duplicate passage name %q", server, n),
				)

				continue
			}

			seen[n] = true

		}
	}

	return errs
}

type SSHServerConfig struct {
	User     string
	Timeout  time.Duration
	Address  string
	Retries  int
	Passages map[string]*PassageConfig
}

const (
	DefaultTimeout = 5 * time.Second
	DefaultRetries = 3
)

func (c *SSHServerConfig) defaults() error {
	if c.User == "" {
		u, err := user.Current()
		if err != nil {
			return err
		}

		c.User = u.Username
	}

	if c.Timeout == 0 {
		c.Timeout = DefaultTimeout
	}

	if c.Retries == 0 {
		c.Retries = DefaultRetries
	}

	return nil
}

func (c *SSHServerConfig) validate(name string) []error {
	if err := c.defaults(); err != nil {
		return []error{err}
	}

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

package server

import (
	"fmt"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type ConfigSuite struct{}

var _ = Suite(&ConfigSuite{})

func (s *ConfigSuite) TestString(c *C) {
	config := &Config{
		Servers: map[string]*SSHServerConfig{
			"foo": {User: "foo", Address: "qux", Passages: map[string]*PassageConfig{
				"qux": {Type: "bar", Address: "foo", Local: "baz"},
			}},
		},
	}

	y, err := config.Marshal()
	c.Assert(err, IsNil)
	fmt.Println(string(y))
}

func (s *ConfigSuite) TestValidate(c *C) {
	config := &Config{
		Servers: map[string]*SSHServerConfig{
			"foo": {User: "foo", Address: "qux", Passages: map[string]*PassageConfig{
				"qux": {Type: "tcp", Address: "foo", Local: "baz"},
			}},
		},
	}

	err := config.Validate()
	c.Assert(err, IsNil)
}

func (s *ConfigSuite) TestValidateErrors(c *C) {
	config := &Config{
		Servers: map[string]*SSHServerConfig{
			"foo": {User: "foo", Address: "", Passages: map[string]*PassageConfig{
				"qux": {Type: "bar", Address: "foo", Local: "baz"},
			}},
			"baz": {},
		},
	}

	err := config.Validate()
	fmt.Println(err)
	c.Assert(err.(*ConfigError).Errors, HasLen, 4)
}

func (s *ConfigSuite) TestValidateEmpty(c *C) {
	config := &Config{}

	err := config.Validate()
	c.Assert(err.Error(), Equals, "invalid empty config")
}

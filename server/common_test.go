package server

import (
	"fmt"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type CommonSuite struct{}

var _ = Suite(&CommonSuite{})

func (s *CommonSuite) AATestString(c *C) {
	config := &Config{
		Servers: []SSHServerConfig{
			{User: "foo", Address: "qux", Passages: []PassageConfig{
				{Remote: RemoteConfig{Type: "bar", Address: "foo"}, Local: "baz"},
			}},
		},
	}

	y, err := config.Marshal()
	c.Assert(err, IsNil)
	fmt.Println(string(y))
}

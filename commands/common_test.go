package commands

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type CommonSuite struct{}

var _ = Suite(&CommonSuite{})

func (s *CommonSuite) TestNewRemote(c *C) {
	r := &Remote{}
	err := r.Set("127.0.0.1:42/tcp")
	c.Assert(err, IsNil)
	c.Assert(r.Remote.String(), Equals, "127.0.0.1:42/tcp")
}

func (s *CommonSuite) TestNewRemotePort(c *C) {
	r := &Remote{}
	err := r.Set("42")
	c.Assert(err, IsNil)
	c.Assert(r.Remote.String(), Equals, "127.0.0.1:42/tcp")

	err = r.Set(":42")
	c.Assert(err, IsNil)
	c.Assert(r.Remote.String(), Equals, "127.0.0.1:42/tcp")
}

func (s *CommonSuite) TestNewRemoteContainer(c *C) {
	r := &Remote{}
	err := r.Set("container=foo:42")
	c.Assert(err, IsNil)
	c.Assert(r.Remote.String(), Equals, "<container=foo>::42/tcp")
}

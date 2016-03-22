package commands

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type CommonSuite struct{}

var _ = Suite(&CommonSuite{})

func (s *CommonSuite) TestRemoteUnmarshalFlag(c *C) {
	r := &Remote{}
	err := r.UnmarshalFlag("127.0.0.1:42/tcp")
	c.Assert(err, IsNil)
	c.Assert(r.Remote.String(), Equals, "127.0.0.1:42/tcp")
}

func (s *CommonSuite) TestRemoteUnmarshalFlagPort(c *C) {
	r := &Remote{}
	err := r.UnmarshalFlag("42")
	c.Assert(err, IsNil)
	c.Assert(r.Remote.String(), Equals, "127.0.0.1:42/tcp")

	r = &Remote{}
	err = r.UnmarshalFlag(":42")
	c.Assert(err, IsNil)
	c.Assert(r.Remote.String(), Equals, "127.0.0.1:42/tcp")
}

func (s *CommonSuite) TestRemoteUnmarshalFlagContainer(c *C) {
	r := &Remote{}
	err := r.UnmarshalFlag("container=foo:42")
	c.Assert(err, IsNil)
	c.Assert(r.Remote.String(), Equals, "<container=foo>::42/tcp")
}

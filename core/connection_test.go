package core

import (
	"golang.org/x/crypto/ssh"
	. "gopkg.in/check.v1"
)

type TunnelSuite struct{}

var _ = Suite(&TunnelSuite{})

func (s *TunnelSuite) TestString(c *C) {
	ssh := NewSSHConnection(
		MustResolveAddr("tcp", "localhost:22"),
		&ssh.ClientConfig{
			User: "root",
		}, 1,
	)

	c.Assert(ssh.String(), Equals, "root@127.0.0.1:22")
}

package passage

import (
	"golang.org/x/crypto/ssh"
	. "gopkg.in/check.v1"
)

type TunnelSuite struct{}

var _ = Suite(&TunnelSuite{})

func (s *TunnelSuite) TestString(c *C) {
	ssh := NewSSHConnection(
		MustResolveTCPAddr("tcp", "localhost:22"),
		&ssh.ClientConfig{
			User: "root",
		},
	)

	c.Assert(ssh.String(), Equals, "root@127.0.0.1:22")
}

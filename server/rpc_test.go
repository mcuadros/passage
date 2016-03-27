package server

import (
	"net"
	"net/rpc"
	"time"

	. "gopkg.in/check.v1"
)

type RPCSuite struct{}

var _ = Suite(&RPCSuite{})

func (s *RPCSuite) TestNewRPC(c *C) {
	config := getConfigFixture()

	server := NewServer()
	err := server.Load(config)
	c.Assert(err, IsNil)
	defer server.Close()

	a, err := net.ResolveUnixAddr("unix", "/tmp/rpc.sock")
	c.Assert(err, IsNil)

	rpcServer := NewRPC(server)
	go rpcServer.Listen(a)
	time.Sleep(100 * time.Millisecond)
	defer rpcServer.Close()

	rpcClient, err := rpc.Dial("unix", rpcServer.l.String())
	c.Assert(err, IsNil)

	var reply string
	err = rpcClient.Call("Server.Addr", "foo", &reply)
	c.Assert(err, IsNil)

	c.Assert(reply, Equals, "[::]:8400")
}

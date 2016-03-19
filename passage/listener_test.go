package passage

import (
	"net"
	"sync"
	"time"

	. "gopkg.in/check.v1"
)

type ListenerSuite struct{}

var _ = Suite(&ListenerSuite{})

func (s *ListenerSuite) TestStart(c *C) {
	local, _ := net.ResolveTCPAddr("tcp", ":0")

	var conn int
	l := NewListener(local)
	l.Handler = func(c net.Conn) error {
		conn++
		return nil
	}

	err := l.Start()
	c.Assert(err, IsNil)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		conn, err := net.Dial(local.Network(), l.String())
		c.Assert(err, IsNil)

		conn.Close()
		time.Sleep(100 * time.Millisecond)
	}()

	wg.Wait()

	err = l.Close()
	c.Assert(err, IsNil)
	c.Assert(conn, Equals, 1)
}

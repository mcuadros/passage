package passage

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"

	. "gopkg.in/check.v1"
)

type RemoteSuite struct{}

var _ = Suite(&RemoteSuite{})

func (s *RemoteSuite) TestNewRemote(c *C) {
	r := NewRemote("tcp", "localhost", "42")
	a, err := r.Addr(nil)
	c.Assert(err, IsNil)
	c.Assert(a.Network(), Equals, "tcp")
	c.Assert(a.String(), Equals, "127.0.0.1:42")
}

func (s *RemoteSuite) TestNewLocalhostRemote(c *C) {
	r := NewLocalhostRemote("tcp", "42")
	a, err := r.Addr(nil)
	c.Assert(err, IsNil)
	c.Assert(a.Network(), Equals, "tcp")
	c.Assert(a.String(), Equals, "127.0.0.1:42")
}

func (s *RemoteSuite) TestRemoteString(c *C) {
	r := NewLocalhostRemote("tcp", "42")
	c.Assert(r.String(), Equals, "localhost:42")
}

func (s *RemoteSuite) TestGetContainerIP(c *C) {
	r := NewContainerRemote("foo", "42")
	a, err := r.Addr(&SSHFixture{})
	c.Assert(err, IsNil)
	c.Assert(a.Network(), Equals, "tcp")
	c.Assert(a.String(), Equals, "172.17.0.2:42")
	c.Assert(r.String(), Equals, "<foo>172.17.0.2:42")
}

func (s *RemoteSuite) TestContainerRemoteString(c *C) {
	r := NewContainerRemote("foo", "42")
	c.Assert(r.String(), Equals, "<foo>::42")
}

type SSHFixture struct{}

func (s *SSHFixture) Conn(a net.Addr) (net.Conn, error) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `[{"Id":"aaa1a054326e8d1bb28617cc04230b16394972d92c9b545c5ab1571f23810524","Names":["/foo"],"Image":"swarm","ImageID":"sha256:291cbe419fe661bfff00d4b2ed7c599f348c7001c17042b2b9b369c495819715","Command":"foo","Created":1458357052,"Ports":[{"PrivatePort":2375,"Type":"tcp"},{"IP":"0.0.0.0","PrivatePort":4000,"PublicPort":4000,"Type":"tcp"}],"SizeRootFs":18106629,"Labels":{},"Status":"Up 39 hours","HostConfig":{"NetworkMode":"default"},"NetworkSettings":{"Networks":{"bridge":{"IPAMConfig":null,"Links":null,"Aliases":null,"NetworkID":"","EndpointID":"c05c5f539e90a400704f8e309cff49b30e53d428217ba2aec28ace6e486851ca","Gateway":"172.17.0.1","IPAddress":"172.17.0.2","IPPrefixLen":16,"IPv6Gateway":"","GlobalIPv6Address":"","GlobalIPv6PrefixLen":0,"MacAddress":"02:42:ac:11:00:02"}}}}]`)
	}))

	url, _ := url.Parse(ts.URL)
	return net.Dial("tcp", url.Host)
}

func (s *SSHFixture) Tunnel(c net.Conn, a net.Addr) error {
	return nil
}

func (s *SSHFixture) String() string {
	return ""
}

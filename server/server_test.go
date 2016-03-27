package server

import . "gopkg.in/check.v1"

type ServerSuite struct{}

var _ = Suite(&ServerSuite{})

func (s *ServerSuite) TestNewServer(c *C) {
	config := getConfigFixture()

	server := NewServer()
	err := server.Load(config)
	c.Assert(err, IsNil)
	defer server.Close()

	c.Assert(server.servers, HasLen, 1)
	c.Assert(server.servers["root@127.0.0.1:22"], NotNil)
	c.Assert(server.passages, HasLen, 3)
	c.Assert(server.passages["foo"], NotNil)
	c.Assert(server.passages["(root@127.0.0.1:22)-[localhost:8400/tcp]"], NotNil)
	c.Assert(server.passages["(root@127.0.0.1:22)-[localhost:8500/tcp]"], NotNil)
}

func (s *ServerSuite) TestLoadChangeServer(c *C) {
	config := getConfigFixture()

	server := NewServer()
	err := server.Load(config)
	c.Assert(err, IsNil)
	defer server.Close()

	c.Assert(server.servers, HasLen, 1)
	c.Assert(server.passages, HasLen, 3)

	config.Servers[0].User = "qux"
	err = server.Load(config)
	c.Assert(err, IsNil)
	c.Assert(server.servers, HasLen, 1)
	c.Assert(server.passages, HasLen, 3)
}

func (s *ServerSuite) TestLoadChangePassage(c *C) {
	config := getConfigFixture()

	server := NewServer()
	err := server.Load(config)
	c.Assert(err, IsNil)
	defer server.Close()

	c.Assert(server.servers, HasLen, 1)
	c.Assert(server.passages, HasLen, 3)

	config.Servers[0].Passages[0].Remote.Type = "container"
	err = server.Load(config)
	c.Assert(err, IsNil)
	c.Assert(server.servers, HasLen, 1)
	c.Assert(server.passages, HasLen, 3)
}

func (s *ServerSuite) TestLoadNoChange(c *C) {
	config := getConfigFixture()

	server := NewServer()
	err := server.Load(config)
	c.Assert(err, IsNil)
	defer server.Close()

	c.Assert(server.servers, HasLen, 1)
	c.Assert(server.passages, HasLen, 3)

	err = server.Load(config)
	c.Assert(err, IsNil)
	c.Assert(server.servers, HasLen, 1)
	c.Assert(server.passages, HasLen, 3)
}

func getConfigFixture() *Config {
	return &Config{
		Servers: []SSHServerConfig{
			{User: "root", Address: "localhost:22", Passages: []PassageConfig{
				{
					Name:   "foo",
					Remote: RemoteConfig{Type: "tcp", Address: "localhost:8400"},
					Local:  ":8400",
				},
				{
					Remote: RemoteConfig{Type: "tcp", Address: "localhost:8400"},
				},
				{
					Remote: RemoteConfig{Type: "tcp", Address: "localhost:8500"},
				},
			}},
		},
	}
}

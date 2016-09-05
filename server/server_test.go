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
	c.Assert(server.servers["baz"], NotNil)
	c.Assert(server.passages, HasLen, 3)
	c.Assert(server.passages["foo"], NotNil)
	c.Assert(server.passages["bar"], NotNil)
	c.Assert(server.passages["qux"], NotNil)
}

func (s *ServerSuite) TestLoadChangeServer(c *C) {
	config := getConfigFixture()

	server := NewServer()
	err := server.Load(config)
	c.Assert(err, IsNil)
	defer server.Close()

	c.Assert(server.servers, HasLen, 1)
	c.Assert(server.servers["baz"].Config().User, Equals, "root")
	c.Assert(server.passages, HasLen, 3)

	config.Servers["baz"].User = "qux"
	err = server.Load(config)
	c.Assert(err, IsNil)
	c.Assert(server.servers, HasLen, 1)
	c.Assert(server.servers["baz"].Config().User, Equals, "qux")
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

	config.Servers["baz"].Passages["foo"].Type = "container"
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
		Servers: map[string]*SSHServerConfig{
			"baz": {
				User:    "root",
				Address: "localhost:22",
				Passages: map[string]*PassageConfig{
					"foo": {
						Type:    "tcp",
						Address: "localhost:8400",
						Local:   ":8400",
					},
					"bar": {
						Type:    "tcp",
						Address: "localhost:8400",
					},
					"qux": {
						Type:    "tcp",
						Address: "localhost:8500",
					},
				}},
		},
	}
}

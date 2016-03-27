package server

import (
	"fmt"
	"net"
	"net/rpc"

	"github.com/mcuadros/passage/core"
)

type RPCServer struct {
	s *Server
	l *core.Listener
	r *rpc.Server
}

func NewRPCServer(s *Server) *RPCServer {
	return &RPCServer{s: s}
}

func (r *RPCServer) Listen(a net.Addr) error {
	r.newRPCServer()
	r.newListener(a)

	return r.l.Start()
}

func (r *RPCServer) newRPCServer() {
	r.r = rpc.NewServer()
	r.r.RegisterName("Server", &RPCContainer{s: r.s})
}

func (r *RPCServer) newListener(a net.Addr) {
	r.l = core.NewListener(a)
	r.l.Handler = func(conn net.Conn) error {
		r.r.ServeConn(conn)
		return nil
	}
}

func (r *RPCServer) Close() error {
	if r.l == nil {
		return nil
	}

	return r.l.Close()
}

type RPCContainer struct {
	s *Server
}

func (r *RPCContainer) Addr(passage string, reply *string) error {
	p, ok := r.s.passages[passage]
	if !ok {
		return fmt.Errorf("unable to find a passage with name %q", passage)
	}

	*reply = p.Addr()
	return nil
}

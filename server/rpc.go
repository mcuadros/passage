package server

import (
	"fmt"
	"net"
	"net/rpc"

	"github.com/mcuadros/passage/core"
)

type RPC struct {
	s *Server
	l *core.Listener
	r *rpc.Server
}

func NewRPC(s *Server) *RPC {
	return &RPC{s: s}
}

func (r *RPC) Listen(a net.Addr) error {
	r.newRPCServer()
	r.newListener(a)

	return r.l.Start()
}

func (r *RPC) newRPCServer() {
	r.r = rpc.NewServer()
	r.r.RegisterName("Server", &RPCContainer{s: r.s})
}

func (r *RPC) newListener(a net.Addr) {
	r.l = core.NewListener(a)
	r.l.Handler = func(conn net.Conn) error {
		r.r.ServeConn(conn)
		return nil
	}
}

func (r *RPC) Close() error {
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

package core

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type ListenerHandler func(net.Conn) error

type Listener struct {
	a net.Addr
	l net.Listener

	Handler     ListenerHandler
	Connections int32
}

func NewListener(a net.Addr) *Listener {
	return &Listener{a: a}
}

func (l *Listener) Start() error {
	var err error
	l.l, err = net.Listen(l.a.Network(), l.a.String())
	if err != nil {
		return fmt.Errorf("error creating listener: %s", err)
	}

	go l.listen()
	return nil
}

func (l *Listener) listen() {
	for {
		conn, err := l.l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		atomic.AddInt32(&l.Connections, 1)
		go func(c net.Conn) {
			err := l.Handler(c)
			if err != nil {
				fmt.Println("error handling connection", err)
			}

			c.Close()
			atomic.AddInt32(&l.Connections, -1)
		}(conn)
	}
}

func (l *Listener) Close() error {
	return l.l.Close()
}

func (l *Listener) String() string {
	if l.l == nil {
		return "<nil>"
	}

	return l.l.Addr().String()
}

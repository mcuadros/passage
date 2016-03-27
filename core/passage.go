package core

import (
	"fmt"
	"net"
)

type Passage struct {
	c SSHConnection
	r Remote
	l *Listener
}

func NewPassage(c SSHConnection, r Remote) *Passage {
	return &Passage{c: c, r: r}
}

func (p *Passage) Start(a net.Addr) error {
	p.buildListener(a)
	return p.l.Start()
}

func (p *Passage) Close() error {
	return p.l.Close()
}

func (p *Passage) buildListener(a net.Addr) {
	p.l = NewListener(a)
	p.l.Handler = func(c net.Conn) error {
		remote, err := p.r.Addr(p.c)
		if err != nil {
			return err
		}

		return p.c.Tunnel(c, remote)
	}
}

func (p *Passage) Addr() string {
	if p.l == nil {
		return "<nil>"
	}

	return p.l.String()
}

func (p *Passage) String() string {
	return fmt.Sprintf("(%s)-[%s]", p.c, p.r)
}

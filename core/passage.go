package core

import "net"

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

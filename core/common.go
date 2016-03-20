package core

import "net"

func MustResolveTCPAddr(network, address string) net.Addr {
	a, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		panic(err)
	}

	return a
}

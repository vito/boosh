package cloudformer

import (
	"net"
)

func CIDR(cidr string) *net.IPNet {
	_, net, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}

	return net
}

func IP(addr string) net.IP {
	ip := net.ParseIP(addr)
	if ip == nil {
		panic("invalid ip: " + addr)
	}

	return ip
}

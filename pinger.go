package main

import (
	"net"
	"strings"

	ping "github.com/digineo/go-ping"
)

func resolveIPAddrAndPinger(addr string) (*net.IPAddr, *ping.Pinger, error) {
	var ra *net.IPAddr
	var pinger *ping.Pinger
	if strings.Contains(addr, ":") {
		r, err := net.ResolveIPAddr("ip6", addr)
		if err != nil {
			return ra, pinger, err
		}
		ra = r
		p, err := ping.New("", "::")
		if err != nil {
			return ra, pinger, err
		}
		pinger = p
	} else {
		r, err := net.ResolveIPAddr("ip4", addr)
		if err != nil {
			return ra, pinger, err
		}
		ra = r
		p, err := ping.New("0.0.0.0", "")
		if err != nil {
			return ra, pinger, err
		}
		pinger = p
	}
	return ra, pinger, nil
}

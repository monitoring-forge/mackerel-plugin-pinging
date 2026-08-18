package main

import (
	"net"

	ping "github.com/digineo/go-ping"
)
...
func resolveIPAddrAndPinger(addr string) (*net.IPAddr, *ping.Pinger, error) {
	r, err := net.ResolveIPAddr("ip", addr)
	if err != nil {
		return nil, nil, err
	}

	var pinger *ping.Pinger
	if r.IP.To4() != nil {
		pinger, err = ping.New("0.0.0.0", "")
	} else {
		pinger, err = ping.New("", "::")
	}
	if err != nil {
		return r, nil, err
	}
	return r, pinger, nil
}
}

package main

import "testing"

// test resolveIPAddrAndPinger function
func TestResolveIPAddrAndPinger(t *testing.T) {
	// Test IPv4 address
	ipv4Addr := "127.0.0.1"
	ipv4IPAddr, ipv4Pinger, err := resolveIPAddrAndPinger(ipv4Addr)
	if err != nil {
		t.Errorf("Failed to resolve IPv4 address: %v", err)
	}
	if ipv4IPAddr == nil || ipv4Pinger == nil {
		t.Errorf("Expected non-nil values for IPv4 address and pinger")
	}

	// Test IPv6 address
	ipv6Addr := "2001:4860:4860::8888"
	ipv6IPAddr, ipv6Pinger, err := resolveIPAddrAndPinger(ipv6Addr)
	if err != nil {
		t.Errorf("Failed to resolve IPv6 address: %v", err)
	}
	if ipv6IPAddr == nil || ipv6Pinger == nil {
		t.Errorf("Expected non-nil values for IPv6 address and pinger")
	}
}

package main

import (
	"errors"
	"os"
	"syscall"
	"testing"
)

// test resolveIPAddrAndPinger function
func TestResolveIPAddrAndPinger(t *testing.T) {
	// Test IPv4 address
	ipv4Addr := "127.0.0.1"
	ipv4IPAddr, ipv4Pinger, err := resolveIPAddrAndPinger(ipv4Addr)
	if err != nil {
		if os.IsPermission(err) || errors.Is(err, syscall.EPERM) || errors.Is(err, syscall.EACCES) {
			t.Skipf("skipping: requires privileges to create ICMP sockets: %v", err)
		}
		t.Fatalf("failed to resolve IPv4 address: %v", err)
	}
	if ipv4IPAddr == nil || ipv4Pinger == nil {
		t.Fatalf("expected non-nil values for IPv4 address and pinger")
	}
	defer ipv4Pinger.Close()

	// Test IPv6 address
	ipv6Addr := "2001:4860:4860::8888"
	ipv6IPAddr, ipv6Pinger, err := resolveIPAddrAndPinger(ipv6Addr)
	if err != nil {
		if os.IsPermission(err) || errors.Is(err, syscall.EPERM) || errors.Is(err, syscall.EACCES) {
			t.Skipf("skipping: requires privileges to create ICMP sockets: %v", err)
		}
		t.Fatalf("failed to resolve IPv6 address: %v", err)
	}
	if ipv6IPAddr == nil || ipv6Pinger == nil {
		t.Fatalf("expected non-nil values for IPv6 address and pinger")
	}
	defer ipv6Pinger.Close()
}

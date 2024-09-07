package main

import (
	"net"
	"testing"
	"time"
)

func setup() {
	arpCache = make(map[string]Host)
	arpEvents = make([]Event, 0)
}

func TestAddHost_InvalidIP(t *testing.T) {
	setup()
	// Test case 2: Invalid IP
	invalidIPStr := "invalid_ip_address"
	err := AddHost(invalidIPStr, "00:1A:2B:3C:4D:60", "Test")
	if err == nil {
		t.Errorf("AddHost with invalid IP %s did not return an error", invalidIPStr)
	}
}

func TestAddHost_InvalidMAC(t *testing.T) {
	setup()
	// Test case 2: Invalid IP
	invalidMACStr := "invalid_mac_address"
	err := AddHost("192.168.7.1", invalidMACStr, "Test")
	if err == nil {
		t.Errorf("AddHost with invalid MAC %s did not return an error", invalidMACStr)
	}
}

func TestAddDeleteHost(t *testing.T) {
	setup()
	if len(arpCache) != 0 {
		t.Errorf("arpCache length is unexpected initial length!\n\tExpected Length: 0\n\tActual Length: %d", len(arpCache))
	}

	// Test case 1: Valid inputs
	ipStr := "192.168.1.3"
	macStr := "00:1A:2B:3C:4D:60"
	comment := "Smartphone"

	expectedIP := net.ParseIP(ipStr)
	if expectedIP == nil {
		t.Fatalf("Invalid test case: could not parse IP %s", ipStr)
	}

	err := AddHost(ipStr, macStr, comment)
	if err != nil {
		t.Errorf("AddHost returned an error: %v", err)
	}

	expectedHost := Host{
		IP:        expectedIP,
		Comment:   comment,
		MAC:       net.HardwareAddr{0x00, 0x1A, 0x2B, 0x3C, 0x4D, 0x60},
		Timestamp: time.Now(),
	}

	// Retrieve the added host from the map
	addedHost, ok := arpCache[ipStr]
	if !ok {
		t.Fatalf("Host with IP %s was not added to the map", ipStr)
	}

	// Ensure we didn't somehow also add a second record
	if len(arpCache) != 1 {
		t.Errorf("arpCache length is unexpected length after add!\n\tExpected Length: 1\n\tActual Length: %d", len(arpCache))
	}

	// Compare the added host with the expected host
	if addedHost.IP.String() != expectedHost.IP.String() ||
		addedHost.Comment != expectedHost.Comment ||
		addedHost.MAC.String() != expectedHost.MAC.String() {
		t.Errorf("AddHost did not add the host correctly. Expected %+v, got %+v", expectedHost, addedHost)
	}

	// Delete the host too
	DeleteHost(ipStr)
	_, ok = arpCache[ipStr]
	if ok {
		t.Errorf("Failed to delete host from arpCache!")
	}
	if len(arpCache) != 0 {
		t.Errorf("arpCache length is unexpected length after delete!\n\tExpected Length: 0\n\tActual Length: %d", len(arpCache))
	}
}

func TestArpEvents(t *testing.T) {
	setup()
	// Assert arpEvents is a 0 length list
	if len(arpEvents) != 0 {
		t.Errorf("arpEvents length is unexpected length!\n\tExpected Length: 0\n\tActual Length: %d", len(arpEvents))
	}

	ipStr := "192.168.1.3"
	macStr := "00:1A:2B:3C:4D:60"
	comment := "Smartphone"

	expectedIP := net.ParseIP(ipStr)
	if expectedIP == nil {
		t.Fatalf("Invalid test case: could not parse IP %s", ipStr)
	}

	// Add the same host twice, this should *not* generate an event.
	err := AddHost(ipStr, macStr, comment)
	if err != nil {
		t.Errorf("AddHost returned an error: %v", err)
	}
	err = AddHost(ipStr, macStr, comment)
	if err != nil {
		t.Errorf("AddHost returned an error: %v", err)
	}
	if len(arpEvents) != 0 {
		t.Errorf("arpEvents length is unexpected length!\n\tExpected Length: 0\n\tActual Length: %d", len(arpEvents))
	}

	// Add the same host again, but with a new MAC, this *should* generate an event
	macStr = "00:1A:2B:3C:4D:61"
	err = AddHost(ipStr, macStr, comment)
	if err != nil {
		t.Errorf("AddHost returned an error: %v", err)
	}
	if len(arpEvents) != 1 {
		t.Errorf("arpEvents length is unexpected length!\n\tExpected Length: 1\n\tActual Length: %d", len(arpEvents))
	}
}

func TestExpireHosts(t *testing.T) {
	setup()
	// Assert arpEvents is a 0 length list
	if len(arpEvents) != 0 {
		t.Errorf("arpEvents length is unexpected length!\n\tExpected Length: 0\n\tActual Length: %d", len(arpEvents))
	}

	ipStr := "192.168.1.3"
	macStr := "00:1A:2B:3C:4D:60"
	comment := "Smartphone"

	expectedIP := net.ParseIP(ipStr)
	if expectedIP == nil {
		t.Fatalf("Invalid test case: could not parse IP %s", ipStr)
	}

	// Add the same host twice, this should *not* generate an event.
	err := AddHost(ipStr, macStr, comment)
	if err != nil {
		t.Errorf("AddHost returned an error: %v", err)
	}

	ExpireHosts(time.Now())
	if len(arpEvents) != 1 {
		t.Errorf("arpEvents length is unexpected length!\n\tExpected Length: 1\n\tActual Length: %d", len(arpEvents))
	}
	if len(arpCache) != 0 {
		t.Errorf("arpCache length is unexpected length after expire!\n\tExpected Length: 0\n\tActual Length: %d", len(arpCache))
	}
}

package main

import (
    "net"
    "testing"
    "time"
)

func TestAddHost(t *testing.T) {
    // Test case 1: Valid inputs
    ipStr := "192.168.1.3"
    macStr := "00:1A:2B:3C:4D:60"
    name := "Smartphone"

    err := AddHost(ipStr, macStr, name)
    if err != nil {
        t.Errorf("AddHost returned an error: %v", err)
    }

    expectedIP := net.ParseIP(ipStr)
    if expectedIP == nil {
        t.Fatalf("Invalid test case: could not parse IP %s", ipStr)
    }

    expectedHost := Host{
        IP:        expectedIP,
        Name:      name,
        MAC:       net.HardwareAddr{0x00, 0x1A, 0x2B, 0x3C, 0x4D, 0x60},
        Timestamp: time.Now(),
    }

    // Retrieve the added host from the map
    addedHost, ok := ipToHostMap[ipStr]
    if !ok {
        t.Fatalf("Host with IP %s was not added to the map", ipStr)
    }

    // Compare the added host with the expected host
    if addedHost.IP.String() != expectedHost.IP.String() ||
        addedHost.Name != expectedHost.Name ||
        addedHost.MAC.String() != expectedHost.MAC.String() {
        t.Errorf("AddHost did not add the host correctly. Expected %+v, got %+v", expectedHost, addedHost)
    }
}

func TestAddHost_InvalidIP(t *testing.T) {
    // Test case 2: Invalid IP
    invalidIPStr := "invalid_ip_address"
    err := AddHost(invalidIPStr, "00:1A:2B:3C:4D:60", "Test")
    if err == nil {
        t.Errorf("AddHost with invalid IP %s did not return an error", invalidIPStr)
    }
}

func TestAddHost_InvalidMAC(t *testing.T) {
    // Test case 2: Invalid IP
    invalidMACStr := "invalid_mac_address"
    err := AddHost("192.168.7.1", invalidMACStr, "Test")
    if err == nil {
        t.Errorf("AddHost with invalid MAC %s did not return an error", invalidMACStr)
    }
}


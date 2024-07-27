package main

import (
	"fmt"
	"net"
	"time"
)

type Host struct {
	IP        net.IP
	Comment   string
	MAC       net.HardwareAddr
	Timestamp time.Time
}

var arpCache = map[string]Host{
	"192.168.1.1": {
		IP:        net.ParseIP("192.168.1.1"),
		Comment:   "Router",
		MAC:       net.HardwareAddr{0x00, 0x1A, 0x2B, 0x3C, 0x4D, 0x5E},
		Timestamp: time.Now(),
	},
	"192.168.1.2": {
		IP:        net.ParseIP("192.168.1.2"),
		Comment:   "Laptop",
		MAC:       net.HardwareAddr{0x00, 0x1A, 0x2B, 0x3C, 0x4D, 0x5F},
		Timestamp: time.Now(),
	},
}

func AddHost(ipStr, macStr, comment string) error {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return fmt.Errorf("invalid IP address: %s", ipStr)
	}

	mac, err := net.ParseMAC(macStr)
	if err != nil {
		return fmt.Errorf("invalid MAC address: %v", err)
	}

	host := Host{
		IP:        ip,
		Comment:   comment,
		MAC:       mac,
		Timestamp: time.Now(),
	}

	arpCache[ipStr] = host
	return nil
}

func DeleteHost(ipStr string) {
	delete(arpCache, ipStr)
}

func main() {
	err := AddHost("192.168.1.3", "00:1A:2B:3C:4D:60", "Smartphone")
	if err != nil {
		fmt.Printf("Error adding host: %v\n", err)
	}

	for ip, host := range arpCache {
		fmt.Printf("IP: %s, Host: %+v\n", ip, host)
	}
}

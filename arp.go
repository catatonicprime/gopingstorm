package main

import (
    "fmt"
    "net"
    "time"
)

type Host struct {
    IP        net.IP
    Name      string
    MAC       net.HardwareAddr
    Timestamp time.Time
}

var ipToHostMap = map[string]Host{
    "192.168.1.1": {
        IP:        net.ParseIP("192.168.1.1"),
        Name:      "Router",
        MAC:       net.HardwareAddr{0x00, 0x1A, 0x2B, 0x3C, 0x4D, 0x5E},
        Timestamp: time.Now(),
    },
    "192.168.1.2": {
        IP:        net.ParseIP("192.168.1.2"),
        Name:      "Laptop",
        MAC:       net.HardwareAddr{0x00, 0x1A, 0x2B, 0x3C, 0x4D, 0x5F},
        Timestamp: time.Now(),
    },
}

func AddHost(ipStr, macStr, name string) error {
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
        Name:      name,
        MAC:       mac,
        Timestamp: time.Now(),
    }

    ipToHostMap[ipStr] = host
    return nil
}

func main() {
    err := AddHost("192.168.1.3", "00:1A:2B:3C:4D:60", "Smartphone")
    if err != nil {
        fmt.Printf("Error adding host: %v\n", err)
    }

    for ip, host := range ipToHostMap {
        fmt.Printf("IP: %s, Host: %+v\n", ip, host)
    }
}


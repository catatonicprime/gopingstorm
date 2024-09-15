package main

import (
	"fmt"
	"net"
)

// getRouteAndInterface returns the local IP address, route interface, the gateway in use, and any errors during discovery
func getRoute(targetIP string) (*net.Interface, *net.UDPAddr, error) {
	// Create a UDP connection to the target IP address to find the route
	conn, err := net.Dial("udp", targetIP+":80")
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()

	// Extract the local address used by this connection
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	// Retrieve all network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, nil, err
	}

	// Match the local IP address with one of the interfaces
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && ipNet.IP.Equal(localAddr.IP) {
				return &iface, localAddr, nil
			}
		}
	}

	return nil, localAddr, fmt.Errorf("interface not found for IP: %s", localAddr.IP)
}

package main

import (
	"fmt"
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
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

func begin_capture() {
	target := "192.168.7.12"
	iface, addr, err := getRoute(target)
	if iface == nil {
		fmt.Printf("Local IP for target %s: %s on interface %s\n", target, addr, "<nil>")
	} else {
		fmt.Printf("Local IP for target %s: %s on interface %s\n", target, addr, iface.Name)
	}

	// Open the network interface for packet capturing
	handle, err := pcap.OpenLive("epair0b", 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Set filter to capture only ICMP packets and ARP packets
	//TODO: append ' or icmp[0] == 0' to process icmp responses too.
	//TODO: compute the expected dstIP of our hsot for ICMP echo replies and filter to only those.
	//TODO: compute the expected dstIP of our host for ARP replies and filter to only those.
	err = handle.SetBPFFilter("(arp and not ether dst ff:ff:ff:ff:ff:ff)")
	if err != nil {
		log.Fatal(err)
	}

	// Start capturing packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	fmt.Println("Capturing packets...")

	// Iterate over captured packets
	for packet := range packetSource.Packets() {
		// Check if this is an ARP packet
		if arpLayer := packet.Layer(layers.LayerTypeARP); arpLayer != nil {
			arp, _ := arpLayer.(*layers.ARP)
			fmt.Printf("ARP Packet:\n")
			fmt.Printf("  Sender Hardware Address: %s\n", net.HardwareAddr(arp.SourceHwAddress))
			fmt.Printf("  Sender Protocol Address: %s\n", net.IP(arp.SourceProtAddress))
			fmt.Printf("  Target Hardware Address: %s\n", net.HardwareAddr(arp.DstHwAddress))
			fmt.Printf("  Target Protocol Address: %s\n", net.IP(arp.DstProtAddress))
			// TODO: maybe do optional rDNS lookups here for the comment.
			ip := net.IP(arp.SourceProtAddress)
			ipStr := ip.String()
			hw := net.HardwareAddr(arp.SourceHwAddress)
			hwStr := hw.String()
			AddHost(ipStr, hwStr, "auto-cached")
		}
		// Check if this is an ICMP packet
		if icmpLayer := packet.Layer(layers.LayerTypeICMPv4); icmpLayer != nil {
			icmp, _ := icmpLayer.(*layers.ICMPv4)
			fmt.Printf("Packet from %s to %s\n", packet.NetworkLayer().NetworkFlow().Src(), packet.NetworkLayer().NetworkFlow().Dst())
			fmt.Printf("  Type: %v Code: %v\n", icmp.TypeCode.Type(), icmp.TypeCode.Code())
			fmt.Printf("  Payload: %v\n", icmp.Payload)
			fmt.Println()
		}
	}
}

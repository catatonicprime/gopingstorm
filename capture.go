package main

import (
	"fmt"
	"log"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func begin_capture(options Options) {
	// TODO: Get all the interfaces we need to listen on, each one would be a capture routine

	for i, cidr := range options.TargetCIDR {
		fmt.Printf("%d, %s\n", i, cidr.String())
	}

	target := "8.8.8.8"
	iface, addr, err := getRoute(target)
	if iface == nil {
		log.Fatal(fmt.Sprintf("Local IP for target %s: %s on interface %s\n", target, addr.String(), "<nil>"))
		return
	}
	fmt.Printf("Local IP for target %s: %s on interface %s\n", target, addr, iface.Name)

	// Open the network interface for packet capturing
	handle, err := pcap.OpenLive(iface.Name, 1600, true, pcap.BlockForever)
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
			//fmt.Printf("ARP Packet:\n")
			//fmt.Printf("  Sender Hardware Address: %s\n", net.HardwareAddr(arp.SourceHwAddress))
			//fmt.Printf("  Sender Protocol Address: %s\n", net.IP(arp.SourceProtAddress))
			//fmt.Printf("  Target Hardware Address: %s\n", net.HardwareAddr(arp.DstHwAddress))
			//fmt.Printf("  Target Protocol Address: %s\n", net.IP(arp.DstProtAddress))
			// TODO: maybe do optional rDNS lookups here for the comment.
			ip := net.IP(arp.SourceProtAddress)
			ipStr := ip.String()
			hw := net.HardwareAddr(arp.SourceHwAddress)
			hwStr := hw.String()
			globalCache.AddHost(ipStr, hwStr, "auto-cached")
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

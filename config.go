package main

import (
	"net"
)

type Route struct {
	// RouteIP is the remote IP we're routing to.
	RouteIP net.IPNet
	// RouteGW is the remote IP of the gateway in use (nil means no gateway).
	RouteGW net.IP
	// RouteAddr is the local IP of the interface used.
	RouteAddr net.IP
	// RouteInterface is the network interface used route.
	RouteInterface net.Interface
}

type Options struct {
	// TargetCIDR is the CIDR to scan.
	TargetCIDR []net.IPNet

	// List of routes used to route to targets.
	Routes []Route

	// InterfaceName an override for the interface to send ICMP/ARP requests on.
	InterfaceName string

	// ArpTimeout is the timeout in nanoseconds to wait for ARP responses.
	ArpTimeout int64

	// IcmpTimeout is the timeout in nanoseconds to wait for ICMP responses.
	IcmpTimeout int64

	// IcmpMaxTTL is the maximum TTL to set on packets when performing traceroute style ICMP echos.
	IcmpMaxTTL int

	// UiRefreshRate is the amount of time to delay between UI refreshes for he Terminal UI (TUI).
	UiRefreshRate int64
}

func ParseArgs() *Options {
	var options Options = make(Options)

	return &options
}

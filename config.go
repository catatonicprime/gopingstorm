package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
	"time"
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

	// InterfaceName is an override for the interface to send ICMP/ARP requests on.
	InterfaceName string

	// ArpTimeout is the timeout in nanoseconds to wait for ARP responses.
	ArpTimeout time.Duration

	// IcmpTimeout is the timeout in nanoseconds to wait for ICMP responses.
	IcmpTimeout time.Duration

	// IcmpMaxTTL is the maximum TTL to set on packets when performing traceroute style ICMP echos.
	IcmpMaxTTL int

	// UiRefreshRate is the amount of time to delay between UI refreshes for the Terminal UI (TUI).
	UiRefreshRate time.Duration
}

// Custom flag for net.IPNet (CIDR parsing)
type ipNetFlag struct {
	ipNet *net.IPNet
}

func (i *ipNetFlag) String() string {
	if i.ipNet != nil {
		return i.ipNet.String()
	}
	return ""
}

func (i *ipNetFlag) Set(value string) error {
	_, ipNet, err := net.ParseCIDR(value)
	if err != nil {
		return fmt.Errorf("invalid CIDR: %s", value)
	}
	i.ipNet = ipNet
	return nil
}

// Custom flag for []net.IPNet (parsing multiple CIDRs)
type ipNetSliceFlag []net.IPNet

func (i *ipNetSliceFlag) String() string {
	var cidrs []string
	for _, cidr := range *i {
		cidrs = append(cidrs, cidr.String())
	}
	return strings.Join(cidrs, ",")
}

func (i *ipNetSliceFlag) Set(value string) error {
	// Split the input string by commas (you can also split by spaces if needed)
	cidrs := strings.Split(value, ",")
	for _, cidr := range cidrs {
		// Trim spaces around each CIDR in case there are spaces after splitting
		cidr = strings.TrimSpace(cidr)
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			return fmt.Errorf("invalid CIDR: %s", cidr)
		}
		*i = append(*i, *ipNet)
	}
	return nil
}

// parseArgs parses the command-line arguments into an Options struct.
// It returns the parsed Options and an error, if any.
func parseArgs() (Options, error) {
	var opts Options

	// Parsing a slice of net.IPNet (CIDRs)
	var targetCIDRs ipNetSliceFlag
	flag.Var(&targetCIDRs, "cidrs", "Target CIDRs to scan (can specify multiple)")

	// Basic string, int, and duration flags
	flag.StringVar(&opts.InterfaceName, "interface", "", "Interface to use for requests")
	flag.DurationVar(&opts.ArpTimeout, "arp-timeout", time.Second*2, "ARP timeout (e.g., 2s)")
	flag.DurationVar(&opts.IcmpTimeout, "icmp-timeout", time.Second*2, "ICMP timeout (e.g., 2s)")
	flag.IntVar(&opts.IcmpMaxTTL, "icmp-max-ttl", 64, "Maximum TTL for ICMP packets")
	flag.DurationVar(&opts.UiRefreshRate, "ui-refresh", time.Second*1, "UI refresh rate")

	// Parse the command-line flags
	flag.Parse()

	// Assign the parsed CIDRs to the options struct
	opts.TargetCIDR = targetCIDRs

	return opts, nil
}

/*
func main() {
	opts, err := parseArgs()
	if err != nil {
		fmt.Println("Error parsing arguments:", err)
		return
	}

	// Print the parsed options
	fmt.Printf("Options: %+v\n", opts)
}
*/

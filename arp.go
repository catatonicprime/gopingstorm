package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"time"
)

// Host represents a host on a network.
type Host struct {
	// IP is the IP address of the host.
	IP net.IP

	// MAC is the MAC address of the host.
	MAC net.HardwareAddr

	// Timestamp is the time the ARP host was observed
	Timestamp time.Time

	// Comment is additional user-defined annotating information for the host.
	Comment string
}

//TODO: Create a lifetime database of hosts observed too.

// arpCache is the current cache of hosts that has been discovered.
var arpCache map[string]Host = make(map[string]Host)

// Event represents ARP events like a host changing MAC. In the future additional events
// such as the initial observeration of a host.
type Event struct {
	// Description describes the Event
	Description string

	// Timestamp is the time the event occurred
	Timestamp time.Time
}

// arpEvents is a list of Events that have occurred.
var arpEvents []Event = make([]Event, 0)

// AddEvent adds an ARP event ot the arpEvents list and logs the event.
func AddEvent(description string, timestamp time.Time) {
	event := Event{
		Description: description,
		Timestamp:   timestamp,
	}
	log.Printf(fmt.Sprintf("%s: %s", timestamp.Format(time.RFC3339), event.Description))
	arpEvents = append(arpEvents, event)
}

// AddHost adds a new host record to the arpCache with a timestamp of the time of addition.
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
		Timestamp: time.Now(), // TODO: assess if this timestamp should be derived from a packet timestamp instead.
	}

	// Log changes to arpCache if necessary
	if currentHost, exists := arpCache[ipStr]; exists && !bytes.Equal(currentHost.MAC, host.MAC) {
		AddEvent(fmt.Sprintf("Host MAC changed from %s to %s\n", currentHost.MAC, host.MAC), host.Timestamp)
	}

	//Update the cache regardless
	arpCache[ipStr] = host
	return nil
}

// ExpireHosts will delete hosts that are older than the time _since_.
func ExpireHosts(since time.Time) {
	for key, host := range arpCache {
		if host.Timestamp.Before(since) {
			AddEvent(fmt.Sprintf("Host %s has expired", key), time.Now())
			delete(arpCache, key)
		}
	}
}

// DeleteHost deletes a specific host
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

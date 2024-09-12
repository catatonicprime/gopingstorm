package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// globalCache is the cache for the whole process
var globalCache = NewARPCache()

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

// TODO: Create a lifetime database of hosts observed too.
type ARPCache struct {
	// RWMutex allows for multiple readers but single writer.
	mutex sync.RWMutex

	// Cache is the current cache of hosts that has been discovered.
	cache map[string]Host
}

func NewARPCache() *ARPCache {
	return &ARPCache{
		cache: make(map[string]Host),
	}
}

func (cache *ARPCache) Length() int {
	cache.mutex.RLock()
	length := len(cache.cache)
	cache.mutex.RUnlock()
	return length
}

// StringToIP returns the net.IP of the IP represented by ipStr or nil and an error.
func StringToIP(ipStr string) (net.IP, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ipStr)
	}
	return ip, nil
}

// StringToMAC returns the net.HardwareAddr of the MAC represented by macStr or nil and an error.
func StringToMAC(macStr string) (net.HardwareAddr, error) {
	mac, err := net.ParseMAC(macStr)
	if err != nil {
		return nil, fmt.Errorf("invalid MAC address: %v", err)
	}
	return mac, nil
}

// AddHost adds a new host record to the cache with a timestamp of the time of addition. Returns list of errors or nil for success.
func (cache *ARPCache) AddHost(ipStr, macStr, comment string) []error {
	var errors []error = make([]error, 0)

	ip, err := StringToIP(ipStr)
	if err != nil {
		errors = append(errors, err)
	}

	mac, err := StringToMAC(macStr)
	if err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return errors
	}

	host := Host{
		IP:        ip,
		Comment:   comment,
		MAC:       mac,
		Timestamp: time.Now(), // TODO: assess if this timestamp should be derived from a packet timestamp instead.
	}

	//Update the cache regardless
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	cache.cache[ipStr] = host
	return nil
}

// Performs an ARP cache lookup and automatically expires the host if needed.
func (cache *ARPCache) Lookup(ipStr string, since time.Time) *Host {
	cache.mutex.RLock()
	host, exists := cache.cache[ipStr]
	cache.mutex.RUnlock()
	if !exists {
		return nil
	}
	// Expire the host if it should be
	if host.Timestamp.Before(since) {
		cache.DeleteHost(ipStr)
		return nil
	}
	return &host
}

// DeleteHost deletes a specific host
func (cache *ARPCache) DeleteHost(ipStr string) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	delete(cache.cache, ipStr)
}

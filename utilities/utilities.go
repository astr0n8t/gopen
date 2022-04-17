package utilities

import (
	"fmt"
	"net/netip"
	"strings"

	"github.com/astr0n8t/gopen/definitions"
)

// Initializes the array of hosts for beginning of evaluation.
// Input: the options struct object, includes addresses, ports, and root.  We only really need addresses.
// Output: ResultStore which is the main memory struct which is used throughout
func InitHosts(variables definitions.Options) definitions.ResultStore {
	store := definitions.ResultStore{}

	ips := processIP(variables.Addresses)

	store.Hosts = make(map[string]definitions.Host, len(ips))

	for _, ip := range ips {
		store.Hosts[ip] = definitions.Host{Address: ip}
	}

	return store
}

// Takes an address expression and decides the following:
// Case a) Address is a single address in which case add it to the array and return
// Case b) Address is a cidr subnet in which case process and return the given array
// Case c) Address is a address range in which case process and return the given array
// Input: single string of either "1.1.1.1", "1.1.1.1/24", "1.1.1.1-1.1.1.2"
// Output: Array of expanded addresses: ["1.1.1.1", "1.1.1.2"]
func processIP(addressExp string) []string {

	var ips []string

	firstIP, err := netip.ParseAddr(strings.Split(strings.Split(addressExp, "/")[0], "-")[0])

	if err != nil {
		panic(fmt.Errorf("unable to parse address, check address syntax"))
	}

	if firstIP.String() == addressExp {
		ips = append(ips, firstIP.String())
	} else {

		if len(strings.Split(addressExp, "/")) > 1 {
			ips = processSubnet(addressExp)
		} else {
			ips = processRange(addressExp)
		}
	}

	return ips
}

// Parses a cidr subnet
// Input: string in format of "1.1.1.0/24"
// Output: Array of expanded addresses: ["1.1.1.1", "1.1.1.2", "1.1.1.3", ...]
func processSubnet(addressExp string) []string {
	var ips []string

	currentAddress, err1 := netip.ParseAddr(strings.Split(addressExp, "/")[0])
	network, err2 := netip.ParsePrefix(addressExp)

	if err1 != nil || err2 != nil {
		panic(fmt.Errorf("unable to parse address CIDR, check address syntax"))
	}

	// Check if this is a single IP address, returns immediately if true
	if network.IsSingleIP() {
		ips = append(ips, currentAddress.String())
		return ips
	}

	// Check if we are at the beginning address, or a valid address
	if network.Contains(currentAddress) && network.Contains(currentAddress.Prev()) {
		ips = append(ips, currentAddress.String())
	}

	if network.Bits() < 8 {
		fmt.Println("Warning: Parsing networks with a netmask less than 8 will use a large amount of memory.")
		fmt.Println("Continuing, but may run out of memory...")
	}

	isValidIP := true
	for isValidIP {
		currentAddress = currentAddress.Next()
		if network.Contains(currentAddress) {
			ips = append(ips, currentAddress.String())
		} else {
			isValidIP = false

			// If we have an invalid address, the last added address is the broadcast
			if len(ips) > 0 {
				ips = ips[:len(ips)-1]
			}
		}
	}

	return ips
}

// Parses an address range
// Input: string in format of "1.1.1.1-1.1.1.2"
// Output: Array of expanded addresses: ["1.1.1.1", "1.1.1.2"]
func processRange(addressExp string) []string {
	var ips []string

	currentAddress, err1 := netip.ParseAddr(strings.Split(addressExp, "-")[0])
	lastAddress, err2 := netip.ParseAddr(strings.Split(addressExp, "-")[1])

	if err1 != nil || err2 != nil {
		panic(fmt.Errorf("unable to parse ending address, check address syntax"))
	}

	ips = append(ips, currentAddress.String())
	for currentAddress != lastAddress {
		currentAddress = currentAddress.Next()
		ips = append(ips, currentAddress.String())
	}

	return ips
}

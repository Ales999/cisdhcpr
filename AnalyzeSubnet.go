// AnalyzeSubnet.go - code from `Rez Moss`:  https://dev.to/rezmoss/in-depth-guide-to-netnetip-prefix-methods-77-4b3c

package main

import (
	"fmt"
	"net/netip"
)

type SubnetInfo struct {
	Network       netip.Prefix
	FirstUsable   netip.Addr
	LastUsable    netip.Addr
	NumHosts      uint64
	BroadcastAddr netip.Addr // IPv4 only
}

func AnalyzeSubnet(prefix netip.Prefix) (SubnetInfo, error) {
	info := SubnetInfo{Network: prefix}

	if !prefix.IsValid() {
		return info, fmt.Errorf("invalid prefix")
	}

	if prefix.Addr().Is4() {
		// IPv4 calculations
		bits := 32 - prefix.Bits()
		info.NumHosts = (1 << bits) - 2 // Subtract network and broadcast

		network := prefix.Addr()
		info.FirstUsable = network.Next()

		// Calculate broadcast address
		broadcast := network
		for i := 0; i < 1<<bits-1; i++ {
			broadcast = broadcast.Next()
		}
		info.BroadcastAddr = broadcast
		info.LastUsable = broadcast.Prev()
	} else {
		// IPv6 calculations
		bits := 128 - prefix.Bits()
		if bits > 64 {
			info.NumHosts = 0 // Too large to represent
		} else {
			info.NumHosts = 1 << bits
		}

		info.FirstUsable = prefix.Addr()
		// IPv6 doesn't use broadcast addresses
		info.LastUsable = info.FirstUsable
		for i := 0; i < 1<<bits-1; i++ {
			info.LastUsable = info.LastUsable.Next()
		}
	}

	return info, nil
}

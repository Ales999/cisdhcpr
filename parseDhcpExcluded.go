package main

import (
	"fmt"
	"net"
	"strings"
)

// Function to parse DHCP excluded IP addresses from cisco line
func parseDhcpExcluded(ipString string) ([]net.IP, error) {
	ips := strings.Split(ipString, " ")
	if len(ips) == 1 {
		return []net.IP{parseSingleIP(ips[0])}, nil
	} else if len(ips) == 2 {
		startIP := parseSingleIP(ips[0])
		endIP := parseSingleIP(ips[1])
		return generateIPRange(startIP, endIP), nil
	}
	return nil, fmt.Errorf("invalid input format")
}

func parseSingleIP(ipString string) net.IP {
	return net.ParseIP(ipString)
}

func generateIPRange(startIP, endIP net.IP) []net.IP {
	start := startIP.To4()
	end := endIP.To4()
	if start == nil || end == nil {
		return nil
	}
	var ipList []net.IP
	for ip := start; !ip.Equal(end) && len(ipList) < 256*256*256; ip = incrementIP(ip) {
		ipList = append(ipList, ip)
	}
	ipList = append(ipList, end)
	return ipList
}

func incrementIP(ipi net.IP) net.IP {
	ip := make(net.IP, len(ipi))
	copy(ip, ipi)
	for i := len(ip) - 1; i >= 0; i-- {
		if ip[i] < 255 {
			ip[i]++
			break
		} else {
			ip[i] = 0
		}
	}
	return ip
}

/*
// Example use:
func main() {
	ips, err := parseDhcpExcluded("192.168.1.1 192.168.1.5")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, ip := range ips {
		fmt.Println(ip.String())
	}
}
*/

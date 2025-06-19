package main

import (
	"bufio"
	"fmt"
	"os"
)

func parseFile(filePath string) (dhcpConfig []DH, err error) {
	// Open the file and read its contents
	file, err := os.Open(filePath)
	if err != nil {
		return []DH{}, err
	}
	defer file.Close()
	// Read the file line
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var textlines []string
	for scanner.Scan() {
		textlines = append(textlines, scanner.Text())
	}
	// Parse the lines and extract DHCP configurations
	dhcpConfig, err = parseCiscoConfig(textlines)
	if err != nil {
		fmt.Println("Error parsing Cisco config:", err)
		return []DH{}, err
	}
	return dhcpConfig, nil
}

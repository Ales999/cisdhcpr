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

	// Вернем только networks, а не hosts, поскольку хосты уже добавлены к типу network.
	var ret []DH
	for _, d := range dhcpConfig {
		// Если это не хост, а сеть, и валидный префикс,  то добавляем в результат.
		if !d.HostType && d.Prefix.IsValid() {
			ret = append(ret, d)
		}
	}
	return ret, nil
}

package main

import (
	"fmt"
	"net/netip"
	"strings"
)

var (
	hostname      string // Имя хоста.
	hostNameFound bool   // Имя хоста в файле найдено или нет.

	foundByIp bool
	//eualip      bool           // Признак что IP совпадают
	ifaceSatus  bool   = true  // Признак что интерфейс не выключен административно (по дефолту - он рабочий)
	secondaryIp bool   = false // Найденный IP это seconary IP
	vrfName     string         // Имя VRF
	faceName    string         // Имя интерфейса
	//onlyip      netip.Addr           // только ip из найденной строки - пример "192.168.1.1"
	gwIp      netip.Addr   // Only Gateway IP address, example: "10.0.0.1"
	netPrefix netip.Prefix // полностью ip адрес и маска - пример: "192.168.1.1/24"
	//aclIn       string               // ACL на IN
	//aclOut      string               // ACL на OUT
)

func parseCiscoConfig(txtlines []string) ([]DH, error) {
	// Parse the lines and return DHCP configurations

	var ipfis []IpFullInfo
	var excludedIPs []netip.Addr // исключенные IP адреса
	var dhcpPools []DH           // список пулов что нашли в данном файле

	// Основной цикл по строкам файла.
	for n, line := range txtlines {
		// Если имя хоста еще не нашли, то проверяем его.
		if !hostNameFound {
			if strings.HasPrefix(line, "hostname") {
				hostNameFound = true
				hostname = line[9:]
				continue
			}
		}
		// Если нашли строку с интерфейсом, то парсим данные по этому интерфейсу.
		if strings.HasPrefix(line, "interface ") && !(strings.Contains(line, "Loopback") || strings.Contains(line, "Tunnel")) {
			ret := parseIpData(&txtlines, line, n)
			if ret != nil {
				ipfis = append(ipfis, ret...)
			}
			continue
		}
		// Если нашли исключение в dhcp
		if strings.HasPrefix(line, "ip dhcp excluded-address") {
			excluded := line[25:] // Берем адрес исключения(ний)
			excluded = strings.TrimSpace(excluded)
			exclIps, err := parseDhcpExcluded(excluded)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// Convert from net.Ip to netip.Addr slice
			for _, ip := range exclIps {
				addr, ok := netip.AddrFromSlice(ip.To4())
				if ok {
					excludedIPs = append(excludedIPs, addr)
				}
			}
		}
		// Если строка содержит dhcp пулл до начинаем обрабатывать данный блок
		if strings.HasPrefix(line, "ip dhcp pool") {
			// Парсим блок данных DHCP настроек
			dhcpPool := parseDhcpPool(&txtlines, line, n)
			dhcpPools = append(dhcpPools, dhcpPool)
		}

	}

	// Test output
	// fmt.Printf("IP Full Info: %v\n", ipfis)
	// fmt.Printf("Excluded IPs: %v\n", excludedIPs)
	// fmt.Printf("DHCP Pools: %v\n", dhcpPools)

	for _, ipf := range ipfis {
		pref := ipf.netPrefix
		if pref.IsValid() {
			//fmt.Printf("IP Prefix is a single IP: %v\n", pref.String())
			for n, pool := range dhcpPools {
				//fmt.Println("Checking DHCP Pools for Gateway:", pool.Scope)
				// Проверяем, совпадает ли адрес шлюза с префиксом IP
				if pool.Gateway.Compare(pref.Addr()) == 0 {
					dhcpPools[n].Vlan = ipf.faceName
				}
			}

		}
	}

	for _, excl := range excludedIPs {
		//fmt.Printf("Excluded IP: %v\n", excl)
		if excl.IsValid() {
			for n, pool := range dhcpPools {
				if pool.Prefix.Contains(excl) {
					dhcpPools[n].Exclusions = append(dhcpPools[n].Exclusions, excl)
					//fmt.Printf("Excluded IP %s added to DHCP Pool: %s\n", excl.String(), pool.Scope)
				}
			}

		}

	}

	// Проверим что GW IP тоже внесен в исключения.
	for i, pool := range dhcpPools {
		if pool.Gateway.IsValid() && pool.Prefix.Contains(pool.Gateway) {
			//pool.Exclusions = append(pool.Exclusions,

			// check curent
			var extAllExcl bool

			for _, excl := range pool.Exclusions {
				if excl.Compare(pool.Gateway) == 0 {
					//fmt.Printf("Gateway IP %s already in exclusions for DHCP Pool: %s\n", pool.Gateway.String(), pool.Scope)
					extAllExcl = true
				}
			}
			// Если адрес шлюза не добавлен в исключения добавим его.
			if !extAllExcl {
				dhcpPools[i].Exclusions = append(dhcpPools[i].Exclusions, pool.Gateway)
				//fmt.Printf("Gateway IP %s added to exclusions for DHCP Pool: %s\n", pool.Gateway.String(), pool.Scope)
			}
		}
	}

	// Переработка DHCP Pools для резервов.
	for _, pool := range dhcpPools {
		// Если это хост то надо найти к какому пулу он относится и добавить его в резервы этого пула.
		if pool.HostType {
			hrs := pool.Reserved
			for _, hr := range hrs {
				if len(hr.HostMac) > 0 {
					//fmt.Println("R-Host IP: ", hr.HostIP.String())
					for j, p := range dhcpPools {
						if !p.HostType {
							if p.Prefix.Contains(hr.HostIP) {
								dhcpPools[j].Reserved = append(dhcpPools[j].Reserved, hr)
								break
							}

						}
					}
				}
			}
		}

	}

	return dhcpPools, nil

}

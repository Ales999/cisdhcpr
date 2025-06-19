package main

import (
	"net/netip"
)

type IpFullInfo struct {
	gwIp      netip.Addr   // IP адрес шлюза.
	netPrefix netip.Prefix // ip адрес c маской  - пример: "192.168.1.1/24"
	hostname  string       // Имя хоста.
	vrfName   string       // Имя VRF
	faceName  string       // Имя интерфейса
	foundByIp bool         // Совпадение найдено
	//eualip      bool         // признак что искомый IP совпадает точно c найденным
	ifaceSatus  bool // Признак что интерфейс не выключен
	secondaryIp bool // Это secondary IP
}

func NewIpFullInfo(
	foundByIp bool,
	//eualip bool,
	ifaceSatus bool,
	secondaryIp bool,
	hostname string,
	vrfName string,
	faceName string,
	gwIp netip.Addr,
	netPrefix netip.Prefix,
) *IpFullInfo {
	return &IpFullInfo{
		foundByIp: foundByIp,
		//eualip:      eualip,
		ifaceSatus:  ifaceSatus,
		secondaryIp: secondaryIp,
		hostname:    hostname,
		vrfName:     vrfName,
		faceName:    faceName,
		gwIp:        gwIp,
		netPrefix:   netPrefix,
	}
}

/*
type AgregInfo struct {
	src []IpFullInfo
	dst []IpFullInfo
}
*/

// String - Перевести в строку данные структуры
func (inf *IpFullInfo) String() {

	//if inf.eualip { // Если искомый ip точно совпадает - выделим цветом и префиксом
	//	fmt.Print("\u001b[31m!>\u001b[32m")
	//}
	var statOff string
	// Если это Secondary IP то выведем это на экран
	if inf.secondaryIp {
		statOff = " (SECNDR)"
	}

	// Если состояние интерфейса как административно выкдюченое (false) - то добавим инфомацию об этом.
	if !inf.ifaceSatus {
		if inf.secondaryIp {
			statOff += "(DOWN)"
		} else {
			statOff += " (DOWN)"
		}
	}

	//fmt.Print("Host: ", inf.hostname, " Iface: ", inf.faceName+statOff, " Vrf: ", inf.vrfName, " IfaceIp: ", inf.netPrefix.String())

	//if inf.eualip {
	//	fmt.Print("\u001b[0m\n")
	//} else {
	//fmt.Print("\n")
	//}

}

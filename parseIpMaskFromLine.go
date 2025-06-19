package main

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
	"strings"
)

// parseIpMaskFromLine - Разбираем строку и возвращаем её IP и Netmask
//
// Input:
// ' ip address 172.24.62.201 255.255.255.248'
//
// Output (by netip.Prefix.String()):
// '172.24.62.201/29'
//
// parseIpMaskFromLine - парсит строку вида ' ip address 172.24.62.201 255.255.255.248'
func parseIpMaskFromLine(line string) (netip.Addr, netip.Prefix, error) {

	// Разбиваем строку на части по пробелам
	cuttingByFour := strings.FieldsFunc(line, func(r rune) bool {
		return r == ' '
	})

	// Проверяем, что строка содержит достаточное количество полей
	if len(cuttingByFour) < 4 {
		return netip.Addr{}, netip.Prefix{}, errors.New("недостаточное количество полей в строке")
	}

	ipStr := cuttingByFour[2]
	ifaceIp, err := netip.ParseAddr(ipStr)
	if err != nil {
		return netip.Addr{}, netip.Prefix{}, fmt.Errorf("ошибка парсинга IP: %w", err)
	}

	parsedMask := cuttingByFour[3]
	stringMask := net.IPMask(net.ParseIP(parsedMask).To4())
	if stringMask == nil {
		return netip.Addr{}, netip.Prefix{}, errors.New("неверная маска подсети")
	}
	lengthMask, _ := stringMask.Size()

	var prefix = netip.PrefixFrom(ifaceIp, lengthMask)

	return ifaceIp, prefix, nil
}

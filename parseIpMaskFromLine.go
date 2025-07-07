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
	parts := strings.FieldsFunc(line, func(r rune) bool {
		return r == ' '
	})

	// Проверяем, что строка содержит достаточное количество полей
	if len(parts) < 4 {
		return netip.Addr{}, netip.Prefix{}, errors.New("недостаточное количество полей в строке")
	}

	ipStr := parts[2]
	addr, err := netip.ParseAddr(ipStr)
	if err != nil {
		return netip.Addr{}, netip.Prefix{}, fmt.Errorf("ошибка парсинга IP: %w", err)
	}

	maskStr := parts[3]
	ipMask := net.IPMask(net.ParseIP(maskStr).To4())
	if ipMask == nil {
		return netip.Addr{}, netip.Prefix{}, errors.New("неверная маска подсети")
	}
	lengthMask, _ := ipMask.Size()

	var prefix = netip.PrefixFrom(addr, lengthMask)

	return addr, prefix, nil
}

package main

import (
	"fmt"
	"log"
	"net/netip"
	"strconv"
	"strings"
)

// parseDhcpPool - парсим данные из строки конфигурации DHCP
func parseDhcpPool(txtlines *[]string, currentLine string, n int) DH {
	dh := DH{} // Создаем пустой объект для заполнения
	reserv := ReservedHost{}
	var err error

	// Получаем имя скопа
	scopeName := strings.TrimSpace(currentLine[12:])
	dh.Scope = scopeName

	startIdx := n + 1
	endIdx := startIdx + 15
	if endIdx > len(*txtlines) {
		endIdx = len(*txtlines)
	}

	for i := startIdx; i < endIdx; i++ {
		line := (*txtlines)[i]
		if !strings.HasPrefix(line, " ") || strings.HasPrefix(line, "!") {
			break
		}
		if after, ok := strings.CutPrefix(line, " network "); ok {
			network := strings.TrimSpace(after)
			fields := strings.Fields(network)
			if len(fields) != 2 {
				log.Fatalf("Ошибка при разборе строки: %v", err)
			}
			// Сохраним для упрощения отладки в переменные
			_onlyIp := fields[0]
			_onlyMask := fields[1]

			addr, err := netip.ParseAddr(_onlyIp)
			if err != nil {
				log.Fatalf("Ошибка при разборе IP адреса: %v", err)
			}
			maskBits, err := ipMaskToBits(_onlyMask)
			if err != nil {
				log.Fatalf("Ошибка при разборе маски: %v", err)
			}
			// Сохраняем длину сети в битах
			dh.MaskBit = maskBits

			prefix := netip.PrefixFrom(addr, maskBits)
			if !prefix.IsValid() {
				log.Fatalln("Неверный префикс сети", network)
			}
			// Сохраняем префикс сети в объекте DHCP
			dh.Prefix = prefix

			subnet, err := AnalyzeSubnet(prefix)
			if err != nil {
				log.Fatalf("Ошибка при анализе подсети: %v", err)
			}
			// Сохраняем адрес стартовый IP и конечный IP подсети который можно выдавать ПК
			dh.StartIP = subnet.FirstUsable
			dh.EndIP = subnet.LastUsable
			// Ищем опции для сети
			for j := i + 1; j < endIdx; j++ {
				nextLine := (*txtlines)[j]
				if !strings.HasPrefix(nextLine, " ") || strings.HasPrefix(nextLine, "!") {
					break
				}
				if strings.Contains(nextLine, " option") {
					dh.Options = append(dh.Options, strings.TrimSpace(nextLine))
				}
			}

		} else if after, ok := strings.CutPrefix(line, " default-router "); ok {
			gwIp := netip.MustParseAddr(strings.TrimSpace(after))
			dh.Gateway = gwIp

		} else if after, ok := strings.CutPrefix(line, " host "); ok {
			parts := strings.Split(strings.TrimSpace(after), " ")
			var hostIp netip.Addr
			if len(parts) >= 2 {
				hostIp = netip.MustParseAddr(parts[0])
			}
			dh.HostType = true
			reserv.HostIP = hostIp
			reserv.HostName = scopeName
			// Ищем даные и опции для хоста
			for j := i + 1; j < endIdx; j++ {
				nextLine := (*txtlines)[j]
				if !strings.HasPrefix(nextLine, " ") || strings.HasPrefix(nextLine, "!") {
					break
				}
				// Перебираем строки блока для хоста
				if after, ok := strings.CutPrefix(nextLine, " hardware-address "); ok {
					mac, err := ParseMacString(after)
					if err != nil {
						fmt.Println("Ошибка при разборе MAC адреса:", err)
					}
					reserv.HostMac = mac
				} else if after, ok := strings.CutPrefix(nextLine, " client-identifier "); ok {
					mac, err := ParseMacString(after)
					if err != nil {
						fmt.Println("Ошибка при разборе MAC адреса:", err)
					}
					reserv.HostMac = mac
				} else if strings.Contains(nextLine, " option") {
					reserv.HostOptions = append(reserv.HostOptions, strings.TrimSpace(nextLine))
				} else if after, ok := strings.CutPrefix(nextLine, " default-router "); ok {
					gwIp := netip.MustParseAddr(strings.TrimSpace(after))
					reserv.GwIP = gwIp
				}

			}

		}
	}

	if len(reserv.HostMac) > 0 {
		dh.Reserved = append(dh.Reserved, reserv)
	}

	return dh
}

func ParseMacString(macstr string) (string, error) {
	macstr = strings.TrimSpace(macstr)
	return macstr, nil
}

// Функция для преобразования IP-маски в количество бит
func ipMaskToBits(mask string) (int, error) {
	// Разделяем строку по точке и преобразуем в целые числа
	parts := strings.Split(mask, ".")
	if len(parts) != 4 {
		return 0, fmt.Errorf("неправильный формат IP-маски")
	}

	var bits int
	for _, part := range parts {
		p, err := strconv.Atoi(part)
		if err != nil || p < 0 || p > 255 {
			return 0, fmt.Errorf("неправильное значение в IP-маске")
		}
		bits += countOnesInBinary(p)
	}
	return bits, nil
}

// Вспомогательная функция для подсчета единиц в двоичном представлении числа
func countOnesInBinary(n int) int {
	count := 0
	for n > 0 {
		if n&1 == 1 {
			count++
		}
		n >>= 1
	}
	return count
}

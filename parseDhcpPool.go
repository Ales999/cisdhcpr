package main

import (
	"fmt"
	"log"
	"net/netip"
	"strconv"
	"strings"
)

func parseDhcpPool(txtlines *[]string, currentLine string, n int) DH {
	var tlsts []string
	dh := DH{} // Создаем пустой объект для заполнения
	reserv := ReservedHost{}

	// Получаем имя скопа
	scopeName := strings.TrimSpace(currentLine[12:])
	dh.Scope = scopeName

	// Если осталось в файле больше 15 строк то берем только 15 строк
	endIndex := n + 15
	if endIndex < len((*txtlines)[n+1:]) {
		tlsts = (*txtlines)[n+1 : endIndex]
	} else { // Иначе - до конца файла
		tlsts = (*txtlines)[n+1:]
	}

	// Бежим по этому блоку
	for _, tlst := range tlsts {
		// Если блок настроек DHCP пула заканчивается то прерываем данный for
		if !strings.HasPrefix(tlst, " ") || strings.HasPrefix(tlst, "!") {
			break
		}
		// Example of line: ' network 192.168.1.0 255.255.255.0'
		if after, ok := strings.CutPrefix(tlst, " network "); ok {
			// * Получаем сеть и маску подсети
			network := strings.TrimSpace(after)

			//network := strings.TrimSpace(strings.Split(tlst, " ")[2])

			// Разбиваем строку на IP и маску подсети
			parts := make([]string, 2)
			n, err := fmt.Sscanf(network, "%s %s", &parts[0], &parts[1])
			if n != 2 || err != nil {
				log.Fatalf("Ошибка при разборе строки: %v", err)
			}

			// Преобразование IP и маски подсети в netip.Addr
			addr, err := netip.ParseAddr(parts[0])
			if err != nil {
				log.Fatalf("Ошибка при разборе адреса: %v", err)
			}
			maskBits, err := ipMaskToBits(parts[1])
			if err != nil {
				log.Fatalf("Ошибка при разборе маски: %v", err)
			}

			// Создание префикса сети
			prefix := netip.PrefixFrom(addr, maskBits)
			// Сохранение префикса сети в структуру DHCP
			if !prefix.IsValid() {
				log.Fatalln("Неверный префикс сети", network)
			}
			dh.Prefix = prefix

			// Теперь получаем начальный и конечный IP диапазона в виде структуры bp AnalyzeSubnet.go
			subnet, err := AnalyzeSubnet(prefix)
			if err != nil {

				log.Fatalf("Ошибка при анализе подсети: %v", err)
			}
			firstIP := subnet.FirstUsable
			lastIP := subnet.LastUsable

			// DEBUG: Вывод начального и конечного IP адресов
			//fmt.Printf("Начальный IP адрес: %s\n", firstIP)
			//fmt.Printf("Конечный IP адрес: %s\n", lastIP)

			dh.StartIP = firstIP
			dh.EndIP = lastIP

		}
		if after, ok := strings.CutPrefix(tlst, " default-router "); ok {
			// * Получаем gateway ip address *
			gwIp := netip.MustParseAddr(strings.TrimSpace(after))
			dh.Gateway = gwIp
		}
		if after, ok := strings.CutPrefix(tlst, " host "); ok {
			// * Получаем host ip address *
			parts := strings.Split(strings.TrimSpace(after), " ")
			var hostIp netip.Addr
			if len(parts) >= 2 {
				hostIp = netip.MustParseAddr(parts[0])
			}
			dh.HostType = true
			reserv.HostIP = hostIp
			reserv.HostName = scopeName
			// Ищем MAC адреса для этого хоста
			// Если осталось в файле больше 15 строк то берем только 15 строк
			endIndex := n + 15
			var mtlsts []string
			if endIndex < len((*txtlines)[n+1:]) {
				mtlsts = (*txtlines)[n+1 : endIndex]
			} else { // Иначе - до конца файла
				mtlsts = (*txtlines)[n+1:]
			}
			for _, mtlst := range mtlsts {
				// Если блок настроек DHCP пула заканчивается то прерываем данный for
				if !strings.HasPrefix(mtlst, " ") || strings.HasPrefix(mtlst, "!") {
					break
				}
				if after, ok := strings.CutPrefix(mtlst, " hardware-address "); ok {
					mac, err := ParseMacString(after)
					if err != nil {
						fmt.Println("Ошибка при разборе MAC адреса:", err)
					}
					reserv.HostMac = mac
				}
				if after, ok := strings.CutPrefix(mtlst, " client-identifier "); ok {
					mac, err := ParseMacString(after)
					if err != nil {
						fmt.Println("Ошибка при разборе MAC адреса:", err)
					}
					reserv.HostMac = mac
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

package main

import (
	"fmt"
	"net/netip"
	"os"
)

type ReservedHost struct {
	HostName string
	HostIP   netip.Addr
	HostMac  string
}

type DH struct {
	Scope      string         // Скоп DHCP
	Vlan       string         // Имя VLAN
	Prefix     netip.Prefix   // Префикс сети
	StartIP    netip.Addr     // Начальный IP-адрес в скоп
	EndIP      netip.Addr     // Конечный IP-адрес в скопе
	Gateway    netip.Addr     // Шлюз
	Exclusions []netip.Addr   // Исключения
	Reserved   []ReservedHost // Резервированные адреса.
	HostType   bool           // Тип хоста false - network, true - host
}

func (d DH) Print() {
	fmt.Println("Имя скопа:", d.Scope)
	fmt.Println("Vlan:", d.Vlan)
	fmt.Println("Начальный IP-адрес:", d.StartIP)
	fmt.Println("Конечный IP-адрес:", d.EndIP)
	fmt.Println("Адрес шлюза:", d.Gateway)
	if len(d.Exclusions) > 0 {
		fmt.Println("Исключения:")
		for _, e := range d.Exclusions {
			fmt.Println("\t", e.String())
		}
	}
	if len(d.Reserved) > 0 {
		fmt.Println("Резервированные IP:")
		for _, r := range d.Reserved {
			if r.HostIP.IsValid() {
				fmt.Println("\t Host:", r.HostName, "IP:", r.HostIP.String(), "MAC:", r.HostMac)
			}
		}
	}
}

func printHostReport(dhcfg []DH) {
	for _, d := range dhcfg {
		if !d.HostType {
			fmt.Println("---")
			d.Print()
		}
	}
	fmt.Println("-----------------------------")
}

/*
Вывод:

1.	Имя скопа
2.	Vlan (Подсеть)
3.	Начальный адрес скопа
4.	Конечный адрес скопа
5.	Адрес шлюза
6.	Исключения (если есть)
7.	Резервированные адреса.
*/

func main() {

	// Получаем аргументы командной строки
	args := os.Args[1:]

	for _, arg := range args {
		filePath, err := GetFullPath(arg)
		if err != nil {
			fmt.Println(err)
		}
		// Проверяем, существует ли файл
		if info, err := os.Stat(filePath); os.IsNotExist(err) || info.IsDir() {
			fmt.Println("Файл не найден или это директория:", arg)
		} else {
			fmt.Println("Файл найден:", arg)
			if ok, err := checkTextFile(filePath); ok && err == nil {
				// Передаем полный путь файла в функцию parseFile
				dhcfg, err := parseFile(arg)
				if err != nil {
					fmt.Println("Ошибка при парсинге файла:", err)
					continue
				}
				// Выводим результат парсинга на экране
				// TODO: Оформить в виде функции вывода на экране
				printHostReport(dhcfg)
				// ---------------
			} else if !ok {
				fmt.Println("Файл не является текстовым")
			}
			if err != nil {
				fmt.Println("Ошибка при проверке файла:", err)
			}
		}
	}
}

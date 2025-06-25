package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/netip"
	"os"
)

type ReservedHost struct {
	HostName    string     `json:"host_name"`              // Имя хоста
	HostIP      netip.Addr `json:"host_ip"`                // IP-адрес х
	GwIP        netip.Addr `json:"gateway_ip"`             // IP-адрес ш
	HostMac     string     `json:"host_mac,omitempty"`     // Адрес MAC хост
	HostOptions []string   `json:"host_options,omitempty"` // Дополнительные опции в DHCP
}

type DH struct {
	Scope      string         `json:"scope"`                // Скоп DHCP
	Vlan       string         `json:"vlan"`                 // Имя VLAN
	Options    []string       `json:"options,omitempty"`    // Дополнительные опции в DHCP
	Prefix     netip.Prefix   `json:"prefix"`               // Префикс сети
	StartIP    netip.Addr     `json:"start_ip"`             // Начальный IP-адрес в скоп
	EndIP      netip.Addr     `json:"end_ip"`               // Конечный IP-адрес в скопе
	Gateway    netip.Addr     `json:"gateway"`              // Шлюз
	Exclusions []netip.Addr   `json:"exclusions,omitempty"` // Исключения
	Reserved   []ReservedHost `json:"reserved,omitempty"`   // Резервированные адреса.
	MaskBit    int            `json:"mask_bit"`             // Скольео бит в сети
	HostType   bool           `json:"-"`                    // Тип хоста false - network, true - host (не выводится в JSON)

}

type Report struct {
	HostName string `json:"host_name"`
	Networks []DH   `json:"networks"`
}

func NewReport(hostname string, networks []DH) *Report {
	return &Report{
		HostName: hostname,
		Networks: networks,
	}
}

func (d DH) Print() {
	fmt.Println("Scope:   ", d.Scope)
	fmt.Println("Vlan:    ", d.Vlan)
	fmt.Println("Mask Bit:", d.MaskBit)
	fmt.Println("Start IP:", d.StartIP)
	fmt.Println("Stop  IP:", d.EndIP)
	fmt.Println("Gateway: ", d.Gateway)
	for _, o := range d.Options {
		fmt.Println("\t", o)
	}
	if len(d.Exclusions) > 0 {
		fmt.Println("Excludes:")
		for _, e := range d.Exclusions {
			fmt.Println("\t", e.String())
		}
	}
	if len(d.Reserved) > 0 {
		fmt.Println("Reserved IP:")
		for _, r := range d.Reserved {
			if r.HostIP.IsValid() {
				fmt.Println("\t Host:", r.HostName, "IP:", r.HostIP.String(), "MAC:", r.HostMac, "GW:", r.GwIP)
				for _, opt := range r.HostOptions {
					fmt.Println("\t\t ", opt)
				}
			}
		}
	}
}

func printHostReport(dhcfg []DH) {
	for _, d := range dhcfg {
		if !d.HostType {
			fmt.Println("-----")
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

	useJson := flag.Bool("json", false, "Вывести в JSON формате")
	// TODO: Output results to file
	// var outputFile string
	// flag.StringVar(&outputFile, "o", "", "Имя файла для вывода")

	flag.Parse()

	/*
		// Получаем аргументы командной строки
		//args := os.Args

		// os.Args содержит все аргументы, включая флаги.
		args := os.Args[1:]
		if len(args) == 0 {
			fmt.Println("Использование: go run main.go [путь к файлу]")
			os.Exit(1)
		}
	*/

	// Получаем оставшиеся аргументы (те, которые не были распознаны как флаги) как имена файлов на обработку
	remainingArgs := flag.Args()
	if len(remainingArgs) > 0 {
		//fmt.Println("Оставшиеся аргументы:", remainingArgs)
		for _, arg := range remainingArgs {
			filePath, err := GetFullPath(arg)
			if err != nil {
				fmt.Println(err)
			}
			// Проверяем, существует ли файл
			if info, err := os.Stat(filePath); os.IsNotExist(err) || info.IsDir() {
				fmt.Println("Файл не найден или это директория:", arg)
			} else {
				// Открываем файл для чтения и проверяем его содержимое
				//fmt.Println("--> Open file:", arg)
				if ok, err := checkTextFile(filePath); ok && err == nil {
					// Передаем полный путь файла в функцию parseFile
					dhcfg, err := parseFile(arg)
					if err != nil {
						fmt.Println("Ошибка при парсинге файла:", err)
						continue
					}
					// Выводим результат парсинга на экране
					// TODO: Оформить в виде функции вывода на экране
					if len(dhcfg) > 0 {
						report := NewReport(arg, dhcfg)
						if *useJson {
							printToJson(report)
						} else {
							fmt.Println("--> Open file:", arg)
							printHostReport(dhcfg)
						}
					}
					// ---------------
				} else if !ok {
					fmt.Println("Файл не является текстовым")
				}
				if err != nil {
					fmt.Println("Ошибка при проверке файла:", err)
				}
			}
		}

	} else {
		fmt.Println("Не указаны файлы для обработки")
		flag.Usage()
	}
}

// Распечатать сообщение в читаемом формате
func printToJson(structToPrint interface{}) error {
	jsonContent, err := json.MarshalIndent(structToPrint, "", "  ")
	if err != nil {
		return err
	}
	// конвертируем []byte в строку и затем печатаем
	fmt.Printf("%+v\n", string(jsonContent))

	return nil
}

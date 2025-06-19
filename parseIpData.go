package main

import "strings"

func parseIpData(txtlines *[]string, currentLine string, n int) (ret []IpFullInfo) {

	faceName = parseInterfaceName(currentLine)

	// Выбираем остатки что еще не сканировали в отдельный слайс (только следующие 20 строк)
	var tlsts []string
	var err error

	// Если осталось в файле больше 22 строк то берем только 21 строку
	endIndex := n + 22
	if endIndex < len((*txtlines)[n+1:]) {
		tlsts = (*txtlines)[n+1 : endIndex]
	} else {
		tlsts = (*txtlines)[n+1:]
	}

	//Очистим от старых записей.
	vrfName = ""
	ifaceSatus = true
	secondaryIp = false

	// Ищем строки с IP и MASK
	for f, tlst := range tlsts {
		// Если блок интерфейса заканчивается то прерываем данный for
		if !strings.HasPrefix(tlst, " ") {
			break
		}
		if strings.HasPrefix(tlst, " vrf forwarding") || strings.HasPrefix(tlst, " ip vrf forwarding") {
			vrfName = parseVrfName(tlst)
		}

		// Если нашли запись об IP/MASK
		if strings.HasPrefix(tlst, " ip address ") {

			gwIp, netPrefix, err = parseIpMaskFromLine(tlst)
			if err != nil {
				continue
			}

			// Если есть! совпадение префикса с искомым, то ищем все остальное.
			//if netPrefix.Contains(findedIp) {
			//	if findedIp.Compare(onlyip) == 0 {
			//		eualip = true
			//	}
			if gwIp.IsValid() && netPrefix.IsValid() {
				foundByIp = true
			}

			secondaryIp = false

			// Проверка что это возможно seconadry
			if strings.Contains(tlst, `secondary`) {
				// Если это так, то установим признак
				secondaryIp = true
			}

			// Создаем новый срез без текущей строки и перебираем его для поиска ACL, если они есть.
			var bodyifaces = tlsts[f+1:]
			for _, body := range bodyifaces {
				// Если обнаружен конец блока (он больше не начинается с пробела) , или `!`, то прекращаем перебор
				if !strings.HasPrefix(body, " ") || strings.HasPrefix(body, "!") {
					break
				}
				// Проверим что интерфейс не выключен
				if strings.Contains(body, "shutdown") && !strings.Contains(body, "description") {
					ifaceSatus = false
				}
				/*
					if strings.HasPrefix(body, " ip access-group") {
						var aclName = parseAclName(body)
						if strings.HasSuffix(body, " in") {
							aclIn = aclName
						}
						if strings.HasSuffix(body, " out") {
							aclOut = aclName
						}
					}
				*/
			} // end for 'bodyifaces'
			//}

		} // End if found needed ip
		if foundByIp {
			ret = append(ret, *NewIpFullInfo(foundByIp /* eualip, */, ifaceSatus, secondaryIp, hostname, vrfName, faceName, gwIp, netPrefix))
		}
		foundByIp = false
		//eualip = false

	} // end for 'tlsts'

	return ret
}

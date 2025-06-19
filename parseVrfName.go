package main

import "strings"

// parseVrfName - Разбираем строку и возвращаем название VRF
func parseVrfName(line string) string {

	// Парсим строку - разложим по частям
	cuttingByTree := strings.FieldsFunc(line, func(r rune) bool {
		return r == ' '
	})

	// Если новый формат
	if strings.HasPrefix(line, " vrf forwarding") {
		return cuttingByTree[2]
	}
	// Иначе старый 'ip vrf forwarding ...'
	return cuttingByTree[3]
}

package main

import "strings"

// parseInterfaceName - Разбираем строку и возвращаем название интерфейса
func parseInterfaceName(line string) string {
	// Парсим строку - разложим по частям
	cuttingByTree := strings.FieldsFunc(line, func(r rune) bool {
		return r == ' '
	})
	return cuttingByTree[1]

}

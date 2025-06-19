package main

import (
	"bufio"
	"os"
	"unicode/utf8"
)

// Проверка того, является ли файл текстовым, проверяя первую строку файла.
func checkTextFile(fileName string) (bool, error) {
	// Открываем файл по указанному имени
	readFile, err := os.Open(fileName)
	if err != nil {
		// Если возникла ошибка при открытии файла, возвращаем false и ошибку
		return false, err
	}
	// Закрываем файл после завершения выполнения функции
	defer readFile.Close()

	// Создаем сканер для чтения файла построчно
	fileScanner := bufio.NewScanner(readFile)
	// Разбиваем файл на строки
	fileScanner.Split(bufio.ScanLines)
	// Читаем первую строку файла
	fileScanner.Scan()

	// Проверяем, является ли первая строка файла корректной UTF-8 строкой
	return utf8.ValidString(string(fileScanner.Text())), nil
}

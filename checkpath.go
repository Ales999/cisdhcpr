package main

import (
	"errors"
	"os"
	"path/filepath"
)

// GetFullPath принимает строку, содержащую имя файла или полный путь,
// и возвращает полный путь к файлу.
func GetFullPath(filePath string) (string, error) {
	if filePath == "" {
		return "", errors.New("пустой аргумент")
	}

	// Если уже полный путь (начинается с `/`),  то возвращаем его как есть.
	if filepath.IsAbs(filePath) {
		return filepath.Clean(filePath), nil
	}

	// Если относительный путь, то соединяем с текущим каталогом.
	currentDir, _ := os.Getwd()
	fullPath := filepath.Join(currentDir, filePath)
	return fullPath, nil
}

/*
// Examlle use:
func main() {
	filePath, err := GetFullPath("example.txt")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(filePath)
	}
}
*/

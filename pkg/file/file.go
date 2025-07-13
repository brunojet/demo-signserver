package file

import (
	"log"
	"os"
	"path/filepath"
)

func CreatePathIfNotExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			panic("Failed to create path: " + err.Error())
		}
		log.Printf("Created path: %s\n", path)
	}
}

func CreateFilePath(filePath string) (*os.File, error) {
	CreatePathIfNotExists(filepath.Dir(filePath))
	return os.Create(filePath)
}

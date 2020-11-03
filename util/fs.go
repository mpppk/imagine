package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func ToRelPath(basePath, path string) (string, error) {
	relP, err := filepath.Rel(basePath, path)
	if err != nil {
		return "", fmt.Errorf("failed to convert path to rel path(%s) from base path(%s): %w", path, basePath, err)
	}
	return filepath.Clean(relP), nil
}

func LoadImagesFromDir(dirPath string, cap int) <-chan string {
	c := make(chan string, cap)
	go func() {
		err := filepath.Walk(dirPath, func(p string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if isImage(p) {
				c <- p
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		close(c)
	}()
	return c
}

func isImage(filePath string) bool {
	imageExts := []string{".png", ".jpg", ".jpeg", ".gif"}
	ext := filepath.Ext(filePath)
	for _, e := range imageExts {
		if ext == e {
			return true
		}
	}
	return false
}

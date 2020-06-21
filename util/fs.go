package util

import (
	"log"
	"os"
	"path/filepath"
)

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

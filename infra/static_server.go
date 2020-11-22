package infra

import (
	"fmt"
	"net/http"
)

func NewFileServer(port uint, basePath string) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.FileServer(http.Dir(basePath)),
	}
}

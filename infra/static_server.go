package infra

import (
	"fmt"
	"net/http"

	"github.com/rakyll/statik/fs"
)

func NewFileServer(port uint, basePath string) *http.Server {
	mux := http.NewServeMux()
	if basePath != "" {
		mux.Handle("/static", http.FileServer(http.Dir(basePath)))
	}
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
}

func NewHtmlServer(port uint) (*http.Server, error) {
	statikFS, err := fs.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize html fs: %w", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(statikFS))
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}, nil
}

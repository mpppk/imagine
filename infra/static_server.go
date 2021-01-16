package infra

import (
	"fmt"
	"net/http"

	_ "github.com/mpppk/imagine/statik"

	"github.com/rakyll/statik/fs"
)

func NewFileServer(port uint, basePath string) *http.Server {
	mux := http.NewServeMux()
	if basePath != "" {
		mux.Handle("/", http.FileServer(http.Dir(basePath)))
	}
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
}

type fileServer struct {
	originalHandler http.Handler
}

func newFileServer(fs http.FileSystem) *fileServer {
	return &fileServer{originalHandler: http.FileServer(fs)}
}

func (f *fileServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Cache-Control", "no-store")
	f.originalHandler.ServeHTTP(res, req)
}

func NewHtmlServer(port uint) (*http.Server, error) {
	statikFS, err := fs.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize html fs: %w", err)
	}
	mux := http.NewServeMux()
	mux.Handle("/", newFileServer(statikFS))
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}, nil
}

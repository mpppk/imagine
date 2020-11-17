package infra

import (
	"fmt"
	"net"
	"net/http"

	"github.com/hydrogen18/stoppableListener"
)

func NewFileServer(port uint, basePath string) (*http.Server, *stoppableListener.StoppableListener, error) {
	server, sl, err := NewStoppableStaticServer(1323)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create stoppable static server: %w", err)
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(basePath))))
	return server, sl, err
}

func NewStoppableStaticServer(port uint) (*http.Server, *stoppableListener.StoppableListener, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, nil, err
	}
	sl, err := stoppableListener.New(l)
	if err != nil {
		return nil, nil, err
	}
	return &http.Server{}, sl, nil
}

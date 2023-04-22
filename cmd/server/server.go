package server

import (
	"fmt"
	"net/http"
)

func New(port uint64) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: nil,
	}
}

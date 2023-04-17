package server

import (
	"context"
	"fmt"
	"ha-video-parser/pkg/service"
	"net"
	"net/http"
)

const keyServerAddr = "serverAddr"

func New(port uint64) *http.Server {

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", service.HelloHandler)

	ctx := context.Background()

	s := &http.Server{
		Addr:    fmt.Sprintf(":%v", port),
		Handler: mux,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}

	return s
}

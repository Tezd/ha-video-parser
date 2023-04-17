package service

import (
	"fmt"
	"io"
	"net/http"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	fmt.Printf("%s: got /hello request\n", ctx.Value("serverAddr"))
	io.WriteString(w, "Hello, HTTP!\n")
}

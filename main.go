package main

import (
	"fmt"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1> Welcome to Woo </h1>")
}

func main() {
	http.HandleFunc("/", index)

	pTLS := NewTLSGenerator("cert.pem", "key.pem")
	pTLS.GenKey()
}

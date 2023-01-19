package main

import (
	"fmt"
	"log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1> Welcome to Woo </h1>")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	serverPort := 3030

	pTLS := NewTLSGenerator("cert.pem", "key.pem")
	pTLS.GenKey()

	// Reference: https://pkg.go.dev/net/http#ListenAndServeTLS
	http.HandleFunc("/", index)
	log.Printf("About to listen on %v. Go to https://127.0.0.1:%v/", serverPort, serverPort)
	err := http.ListenAndServeTLS(fmt.Sprintf(":%v", serverPort), pTLS.certFile, pTLS.keyFile, nil)
	log.Fatal(err)
}

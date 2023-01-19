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
	certFile := "cert.pem"
	keyFile := "key.pem"

	pTLS := NewTLSGenerator(certFile, keyFile)
	pMaster := NewMaster()

	// Reference: https://pkg.go.dev/net/http#ListenAndServeTLS
	http.HandleFunc("/", index)
	http.HandleFunc("/submit", pMaster.MasterTaskSubmission)
	log.Printf("About to listen on %v. Go to https://127.0.0.1:%v/", serverPort, serverPort)
	err := http.ListenAndServeTLS(fmt.Sprintf(":%v", serverPort), pTLS.certFile, pTLS.keyFile, nil)
	log.Fatal(err)
}

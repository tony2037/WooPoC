// Reference: https://go.dev/src/crypto/tls/generate_cert.go

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

type TLSGenerator struct {
	init     bool
	certFile string
	keyFile  string
}

func NewTLSGenerator(_certFile string, _keyFile string) *TLSGenerator {
	fmt.Println("Key Generator")
	instance := new(TLSGenerator)
	instance.certFile = _certFile
	instance.keyFile = _keyFile
	instance.IsKeyExists()
	if !instance.init {
		fmt.Printf("Did not find %v &/ %v \n", instance.certFile, instance.keyFile)
		instance.GenKey()
	}
	return instance
}

func (pInstance *TLSGenerator) IsKeyExists() bool {
	_, certFileErr := os.Stat(filepath.Join("./", pInstance.certFile))
	_, keyFileErr := os.Stat(filepath.Join("./", pInstance.keyFile))
	pInstance.init = !(os.IsNotExist(certFileErr) || os.IsNotExist(keyFileErr))
	return pInstance.init
}

func (pInstance *TLSGenerator) GenKey() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Generate key ...")

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	// ECDSA, ED25519 and RSA subject keys should have the DigitalSignature
	// KeyUsage bits set in the x509.Certificate template
	keyUsage := x509.KeyUsageDigitalSignature
	// Only RSA subject keys should have the KeyEncipherment KeyUsage bits set. In
	// the context of TLS this KeyUsage is particular to RSA key exchange and
	// authentication.
	keyUsage |= x509.KeyUsageKeyEncipherment

	var notBefore time.Time
	notBefore = time.Now()
	notAfter := notBefore.Add(time.Hour * 24 * 365) // Vaild for a year

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed to generate serial number: %v", err)
	}
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Ztex Co"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}
	certOut, err := os.Create(pInstance.certFile)
	if err != nil {
		log.Fatalf("Failed to open %v for writing: %v", pInstance.certFile, err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		log.Fatalf("Failed to write data to cert.pem: %v", err)
	}
	if err := certOut.Close(); err != nil {
		log.Fatalf("Error closing cert.pem: %v", err)
	}
	log.Print("wrote cert.pem\n")

	keyOut, err := os.OpenFile(pInstance.keyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Failed to open %v for writing: %v", pInstance.keyFile, err)
		return
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		log.Fatalf("Unable to marshal private key: %v", err)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		log.Fatalf("Failed to write data to key.pem: %v", err)
	}
	if err := keyOut.Close(); err != nil {
		log.Fatalf("Error closing key.pem: %v", err)
	}
	log.Print("wrote key.pem\n")
}

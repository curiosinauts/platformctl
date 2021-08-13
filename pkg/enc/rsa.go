package enc

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
)

// CREDIT: https://www.systutorials.com/how-to-generate-rsa-private-and-public-key-pair-in-go-lang/

func GenerateRSAKeys() (string, string) {
	// generate key
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
	}
	publickey := &privatekey.PublicKey

	var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(privatekey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	buf := make([]byte, 0, 1024)
	privatePem := bytes.NewBuffer(buf)
	err = pem.Encode(privatePem, privateKeyBlock)
	if err != nil {
		log.Fatalf("error when encoding private pem: %s \n", err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publickey)
	if err != nil {
		log.Fatalf("error when dumping publickey: %s \n", err)
	}
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	if err != nil {
		log.Fatalf("error when create public.pem: %s \n", err)
	}

	buf2 := make([]byte, 0, 1024)
	publicPem := bytes.NewBuffer(buf2)
	err = pem.Encode(publicPem, publicKeyBlock)
	if err != nil {
		log.Fatalf("error when encode public pem: %s \n", err)
	}

	return privatePem.String(), publicPem.String()
}

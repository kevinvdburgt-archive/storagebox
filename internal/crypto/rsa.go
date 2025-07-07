package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
)

func ParseRSAPublicKey(pemData []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pemData)
	if block == nil {
		log.Fatal("Invalid public key PEM")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	return pub.(*rsa.PublicKey)
}

func ParseRSAPrivateKey(pemData []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(pemData)
	if block == nil {
		log.Fatal("Invalid private key PEM")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	return priv
}

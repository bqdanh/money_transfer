package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func GetPrivateKey(path string) *rsa.PrivateKey {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	block, _ := pem.Decode(bytes)
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	return key
}

func GetPublicKey(path string) *rsa.PublicKey {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	block, _ := pem.Decode(bytes)
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	pkey, ok := key.(*rsa.PublicKey)
	if !ok {
		panic("not public key")
	}
	return pkey
}

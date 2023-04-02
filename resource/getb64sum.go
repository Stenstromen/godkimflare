package resource

import (
	"crypto/rsa"
	"crypto/x509"

	b64 "encoding/base64"

	"encoding/pem"
	"log"
	"os"
)

func getb64sum(certfile string) string {
	priv, err := os.ReadFile(certfile)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	block, _ := pem.Decode([]byte(priv))

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	publicKeyDer, err := x509.MarshalPKIXPublicKey(&key.(*rsa.PrivateKey).PublicKey)
	if err != nil {
		log.Fatal(err)
	}

	uEnc := b64.StdEncoding.EncodeToString(publicKeyDer)
	return string(uEnc)
}

package storage

import (
	"compress/gzip"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/devedbv/storagebox/internal/crypto"
	"github.com/devedbv/storagebox/internal/env"
)

func Download(remotePath string) {
	baseURL := env.Get("BASE_URL")
	username := env.Get("USERNAME")
	password := env.Get("PASSWORD")
	privatePem := env.Get("PRIVATE_KEY")
	privateKey := crypto.ParseRSAPrivateKey([]byte(privatePem))

	req, _ := http.NewRequest("GET", baseURL+"/"+remotePath, nil)
	req.SetBasicAuth(username, password)
	r, err := http.DefaultClient.Do(req)
	if err != nil || r.StatusCode != 200 {
		log.Fatalf("Download failed: %v", err)
	}
	defer r.Body.Close()

	prefix := make([]byte, 2)
	io.ReadFull(r.Body, prefix)
	keyLen := int(prefix[0])<<8 | int(prefix[1])
	encryptionKey := make([]byte, keyLen)
	io.ReadFull(r.Body, encryptionKey)

	aesKey, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptionKey)
	if err != nil {
		log.Fatalf("RSA decryption failed: %v", err)
	}

	stream, _ := crypto.NewAESStream(aesKey)
	reader := &cipher.StreamReader{S: stream, R: r.Body}
	gz, _ := gzip.NewReader(reader)
	defer gz.Close()
	io.Copy(os.Stdout, gz)
}

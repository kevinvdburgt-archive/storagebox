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

func Upload(remotePath string) {
	baseURL := env.Get("BASE_URL")
	username := env.Get("USERNAME")
	password := env.Get("PASSWORD")
	publicPem := env.Get("PUBLIC_KEY")

	publicKey := crypto.ParseRSAPublicKey([]byte(publicPem))
	err := mkdirs(baseURL, username, password, remotePath)
	if err != nil {
		log.Fatalf("Failed to create directories: %v", err)
	}

	reader, writer := io.Pipe()

	go func() {
		defer writer.Close()

		aesKey := make([]byte, 32)
		rand.Read(aesKey)
		encryptionKey, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, aesKey)
		if err != nil {
			writer.CloseWithError(err)
			return
		}

		stream, err := crypto.NewAESStream(aesKey)
		if err != nil {
			writer.CloseWithError(err)
			return
		}

		cipherWriter := &cipher.StreamWriter{S: stream, W: writer}
		prefix := []byte{byte(len(encryptionKey) >> 8), byte(len(encryptionKey))}
		writer.Write(prefix)
		writer.Write(encryptionKey)

		gz := gzip.NewWriter(cipherWriter)
		defer gz.Close()
		io.Copy(gz, os.Stdin)
	}()

	req, _ := http.NewRequest("PUT", baseURL+"/"+remotePath, reader)
	req.SetBasicAuth(username, password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil || (resp.StatusCode != 200 && resp.StatusCode != 201) {
		log.Fatalf("Upload failed: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("File uploaded to %s/%s", baseURL, remotePath)
}

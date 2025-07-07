package main

import (
	"fmt"
	"log"

	"os"

	"github.com/devedbv/storagebox/internal/storage"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	if len(os.Args) != 3 {
		fmt.Println("Usage: storagebox <upload|download> <remote-path>")
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "upload":
		storage.Upload(os.Args[2])
	case "download":
		storage.Download(os.Args[2])
	default:
		log.Fatalf("Unknown command: %s", command)
	}
}

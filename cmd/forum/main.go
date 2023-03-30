package main

import (
	// "context"
	"fmt"

	"github.com/Andreifx02/forum/internal/server"
	postrgres "github.com/Andreifx02/forum/internal/storage/postgres"
)

func main() {
	storage, err := postrgres.NewStorage("localhost", 5432, "postgres", "postgrespw", "postgres")
	if err != nil {
		fmt.Printf("Can not connect db: %s", err.Error())
	}

	server := server.New(storage)

	server.StartListen("localhost", 8080)
}

package main

import (
	// "context"
	"log"
	"sync"

	"github.com/Andreifx02/forum/internal/config"
	"github.com/Andreifx02/forum/internal/server"
	postrgres "github.com/Andreifx02/forum/internal/storage/postgres"

)

func main() {
	cfg := config.GetConfig()
	
	storage, err := postrgres.NewStorage(cfg)
	if err != nil {
		log.Printf("Can not connect db: %s\n", err.Error())
	} else {
		log.Println("Connected to postrges")
	}
	
	server := server.New(storage)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		log.Println("Start listen server")
		server.StartListen(cfg)
		wg.Done()
	}()
	
	wg.Wait()
}

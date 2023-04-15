package main

import (
	//"context"
	//"encoding/json"
	//"fmt"
	"log"
	"sync"

	//"time"

	"github.com/Andreifx02/forum/internal/bot"
	"github.com/Andreifx02/forum/internal/config"
	//"github.com/Andreifx02/forum/internal/domain"
	"github.com/Andreifx02/forum/internal/server"
	postrgres "github.com/Andreifx02/forum/internal/storage/postgres"
)

func main() {
	cfg := config.GetConfig()
	
	storage, err := postrgres.NewStorage(cfg)
	if err != nil {
		log.Fatalf("Can not connect db: %s\n", err.Error())
	} else {
		log.Println("Connected to postrges")
	}
	
	server := server.New(storage)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		log.Println("Start listen server")
		err := server.StartListen(cfg)
		if err != nil {
			log.Fatalf("Can not listen: %s\n", err.Error())	
		}
		wg.Done()
	}()

	bot, err := bot.New(cfg)
	if err != nil {
		log.Fatalf("Can not start bot: %s\n", err.Error())
	}
	wg.Add(1)
	go func() {
		log.Println("Start bot")
		bot.Run()
		wg.Done()
	}()

	wg.Wait()
}

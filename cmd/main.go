package main

import (
	"fmt"
	"log"
	"matterpoll-bot/config"
	"matterpoll-bot/internal/entities"
	"matterpoll-bot/internal/handlers"
	"matterpoll-bot/internal/services"
	"matterpoll-bot/internal/storage"
	"matterpoll-bot/internal/storage/database"
	"matterpoll-bot/internal/storage/memory"
	"net/http"
	"time"

	"github.com/mattermost/mattermost-server/v6/model"
)

func main() {
	var store storage.StoreInterface

	entities.Bot = model.NewAPIv4Client(config.ServerURL)
	entities.Bot.SetToken(config.BotToken)
	
	if err := services.RegisterCommands(); err != nil {
		log.Fatalln(err)
	}

	switch config.Mode {
	case "memory":
		store = memory.NewMemoryStore()
		log.Println("Using memory store")
	case "database":
		conn, err := database.NewDatabaseConection()
		if err != nil {
			log.Fatal(err)
		}

		store = database.NewDatabaseStore(conn)
		log.Println("Using database store")
	default:
		log.Fatalf("config.Mode is empty in /internal/config/config.go")
	}

	pollService := services.NewPollService(store)
	mux := http.NewServeMux()

	mux.HandleFunc("/poll-create", handlers.CreatePoll(pollService))
	mux.HandleFunc("/poll-vote", handlers.Vote(pollService))
	mux.HandleFunc("/poll-results", handlers.GetPollResults(pollService))
	mux.HandleFunc("/poll-close", handlers.ClosePoll(pollService))
	mux.HandleFunc("/poll-delete", handlers.DeletePoll(pollService))

	serv := &http.Server{
		Addr:         config.BotSocket,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Println("Bot is running ...")
	if err := serv.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}

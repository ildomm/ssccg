package main

import (
	"context"
	"errors"
	"github.com/ildomm/ssccg/api"
	"github.com/ildomm/ssccg/dao"
	"github.com/ildomm/ssccg/persistence"
	"github.com/ildomm/ssccg/system"
	"log"
	"net/http"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize log standards
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("[SSCCG] ")

	// Initialize database
	// Dev note: this could be replaced with a real database, like Postgres
	querier, err := persistence.NewInMemoryQuerier(ctx)
	if err != nil {
		log.Fatal("Could not initialize database")
	}

	// Initialize services
	deviceDAO := dao.NewDeviceDAO(querier)

	// Initialize the server
	server := api.NewServer()
	if listenAddress := system.ExtractServerPort(); listenAddress != nil {
		server.WithListenAddress(*listenAddress)
	}
	server.WithDeviceManager(deviceDAO)
	log.Println("Starting server on", server.ListenAddress())

	if err := server.Run(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Could not start server on", server.ListenAddress())
		} else {
			log.Println("Server closed")
		}
	}
}

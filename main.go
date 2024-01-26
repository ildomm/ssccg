package main

import (
	"errors"
	"github.com/ildomm/ssccg/api"
	"github.com/ildomm/ssccg/system"
	"log"
	"net/http"
)

func main() {
	// Initialize log standards
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("[SSCCG] ")

	// TODO: Initialize database
	// TODO: Initialize services

	// Initialize the server
	server := api.NewServer()
	if listenAddress := system.ExtractServerPort(); listenAddress != nil {
		server.WithListenAddress(*listenAddress)
	}
	log.Println("Starting server on ", server.ListenAddress())

	if err := server.Run(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Could not start server on ", server.ListenAddress())
		} else {
			log.Println("Server closed")
		}
	}
}

// TODO: Initialize signers list
//
//nolint:all
func cryptoSigners() {

}

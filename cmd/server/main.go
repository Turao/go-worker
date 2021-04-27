package main

import (
	"log"

	"github.com/turao/go-worker/server"
)

func main() {
	server := server.NewServer(":8080")
	log.Println("Serving the new API server...")
	server.ListenAndServe()
}

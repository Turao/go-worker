package main

import (
	"log"

	"github.com/turao/go-worker/pkg/server"
)

func main() {
	server := server.New(":8080")
	log.Println("Serving the new API server...")
	server.ListenAndServe()
}

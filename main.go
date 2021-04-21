package main

import (
	"log"

	"github.com/turao/kami-go/apiserver"
)

func main() {
	server := apiserver.NewServer(":8080")
	log.Println("Serving the new API server...")
	server.ListenAndServe()
}

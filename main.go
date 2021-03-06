package main

import (
	"log"
	"net/url"
	"time"

	"github.com/turao/go-worker/pkg/client"
	"github.com/turao/go-worker/pkg/server"
)

func main() {
	server := server.New(":8080")
	log.Println("Serving the new API server...")
	go func() {
		server.ListenAndServe()
	}()

	log.Println("waiting for server to start")
	time.Sleep(2 * time.Second)
	log.Println("server should have started by now...")

	url, err := url.Parse("http://localhost:8080/job")
	if err != nil {
		log.Fatalln("failed to parse client url")
	}
	c := client.New(url)
	log.Println("sending command via client")
	res, err := c.Dispatch("ls", "-lah")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(res)

	time.Sleep(5 * time.Second)
}

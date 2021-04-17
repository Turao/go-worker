package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/turao/kami-go/apiserver"
)

func main() {
	server := apiserver.NewServer()
	ctx := context.Background()

	jobId, err := server.Service.Dispatch(ctx, "ls", "-lah")
	if err != nil {
		log.Fatalln(err.Error())
	}

	// time.Sleep(1 * time.Second)
	// worker.Stop(jobId)
	time.Sleep(3 * time.Second)

	info, err := server.Service.QueryInfo(ctx, jobId)
	if err != nil {
		log.Fatalln(err.Error())
	}

	d, _ := json.MarshalIndent(info, "", "  ")
	log.Println(string(d))

	time.Sleep(5 * time.Second)

	logs, err := server.Service.QueryLogs(ctx, jobId)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("output:\n", logs.Output)
	log.Println("errors:\n", logs.Errors)

	log.Println("Serving the new API server...")
	server.Serve()
}

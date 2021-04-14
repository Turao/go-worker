package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/turao/kami-go/worker"
)

func main() {
	worker := worker.NewWorker()

	jobId, err := worker.Dispatch("ls", "-lah")
	if err != nil {
		log.Fatalln(err.Error())
	}

	// time.Sleep(1 * time.Second)
	// worker.Stop(jobId)
	time.Sleep(3 * time.Second)

	info, err := worker.QueryInfo(jobId)
	if err != nil {
		log.Fatalln(err.Error())
	}

	d, _ := json.MarshalIndent(info, "", "  ")
	log.Println(string(d))

	time.Sleep(5 * time.Second)

	logs, err := worker.QueryLogs(jobId)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("output:\n", logs.Output)
	log.Println("errors:\n", logs.Errors)
}

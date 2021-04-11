package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/turao/kami-go/worker"
)

func main() {
	worker := worker.MakeWorker()

	jobId, err := worker.Dispatch("ls")
	if err != nil {
		log.Fatalln(err.Error())
	}

	time.Sleep(1 * time.Second)
	worker.Stop(jobId)
	time.Sleep(10 * time.Second)

	info, err := worker.QueryInfo(jobId)
	if err != nil {
		log.Fatalln(err.Error())
	}

	d, _ := json.MarshalIndent(info, "", "  ")
	log.Println(string(d))

	time.Sleep(10 * time.Second)
}

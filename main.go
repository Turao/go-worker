package main

import (
	"encoding/json"
	"log"
	"os/exec"
	"time"

	"github.com/turao/kami-go/worker"
)

func main() {
	worker := worker.MakeWorker()

	command := exec.Command("sleep", "2")

	jobId, err := worker.Dispatch(command)
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

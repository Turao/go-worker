package main

import (
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

	info, err := worker.QueryInfo(jobId)
	if err != nil {
		log.Fatalln(err.Error())
	}

	time.Sleep(1 * time.Second)
	worker.Stop(jobId)

	log.Println(info)

	time.Sleep(10 * time.Second)
}

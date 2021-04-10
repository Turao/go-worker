package main

import (
	"log"
	"os/exec"

	"github.com/turao/kami-go/worker"
)

func main() {
	worker := worker.MakeWorker()

	command := exec.Cmd{}

	jobId, err := worker.Dispatch(command)
	if err != nil {
		log.Fatalln(err.Error())
	}

	worker.Stop(jobId)

	info, err := worker.QueryInfo(jobId)
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println(info)
}

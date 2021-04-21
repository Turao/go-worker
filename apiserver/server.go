package apiserver

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type apiserver struct {
	server *http.Server
	Service
}

func NewServer() *apiserver {
	service := NewWorkerService()
	service = loggingMiddleware{next: service}

	return &apiserver{
		server: &http.Server{
			Addr:    ":8080",
			Handler: makeHandler(service),
		},
		Service: service,
	}
}

func (s *apiserver) Serve() {
	log.Println("[server]", "called")
	errs := make(chan error, 1)
	go func() {
		log.Println("[server]", "serving...")
		errs <- s.server.ListenAndServe()
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("[server] interrupted: %s", <-c)
	}()

	log.Println("[server]", "terminated", <-errs)
}

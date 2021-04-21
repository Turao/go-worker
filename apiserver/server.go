package apiserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type apiserver struct {
	server *http.Server
}

func NewServer(addr string) *apiserver {
	workerservice := newWorkerService()
	workerservice = loggingMiddleware{next: workerservice}

	return &apiserver{
		server: &http.Server{
			Addr:    addr,
			Handler: makeHandler(workerservice),
		},
	}
}

func (s *apiserver) Close() error {
	return s.server.Close()
}

func (s *apiserver) ListenAndServe() {
	log.Println("[server]", "listen and serve")
	errs := make(chan error, 1)

	go func(errs chan<- error) {
		log.Println("[server]", "serving on", s.server.Addr)
		errs <- s.server.ListenAndServe()
	}(errs)

	go func(errs chan<- error) {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("[server] interrupted: %s", <-c)
	}(errs)

	<-errs // blocks until listen throws error or interrupted

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := s.server.Shutdown(ctx)
	if err != nil {
		log.Fatalln("failed to shutdown server")
	}

	log.Println("[server]", "server shutdown")
}

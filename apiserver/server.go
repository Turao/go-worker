package apiserver

type server struct {
	Service Service // todo: stop exposing the server's Service (doing it so I can test in main() now)
}

func NewServer() *server {
	service := NewWorkerService()
	service = loggingMiddleware{next: service}

	return &server{
		Service: service,
	}
}

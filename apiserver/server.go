package apiserver

type server struct {
	service Service
}

func NewServer() {
	service := NewWorkerService()
	service = loggingMiddleware{next: service}
}

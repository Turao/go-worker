package v1

type DispatchRequest struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}

type DispatchResponse struct {
	ID JobID `json:"id"`
}

type StopRequest struct {
	ID JobID `json:"id"`
}

type StopResponse struct {
}

type QueryInfoRequest struct {
	ID JobID `json:"id"`
}

type QueryInfoResponse struct {
	JobInfo // embbeded
}

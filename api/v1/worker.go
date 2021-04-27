package v1

type DispatchRequest struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}

type DispatchResponse struct {
	ID string `json:"id"`
}

type StopRequest struct {
	ID string `json:"id"`
}

type StopResponse struct {
}

type QueryInfoRequest struct {
	ID string `json:"id"`
}

type QueryInfoResponse struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	ExitCode int    `json:"exitCode"`
	Output   string `json:"output"`
	Errors   string `json:"errors"`
}

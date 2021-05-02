package v1

type JobID string

type JobInfo struct {
	ID       JobID  `json:"id"`
	Status   string `json:"status"`
	ExitCode int    `json:"exitCode"`
	Output   string `json:"output"`
	Errors   string `json:"errors"`
}

package redis

const (
	StatusProcessing = "Processing..."
	StatusError      = "Error"
	StatusSuccess    = "Success"
)

type TaskError struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Retries   int    `json:"retries"`
}

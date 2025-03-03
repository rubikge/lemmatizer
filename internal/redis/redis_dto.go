package redis

type Response struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

const (
	StatusProcessing = "Processing..."
	StatusDone       = "Done"
	StatusError      = "Error"
)

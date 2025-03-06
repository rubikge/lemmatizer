package dto

type SearchResult struct {
	Status       string  `json:"status"`
	ProductTitle string  `json:"product_title"`
	TotalScore   float64 `json:"total_score"`
	TaskID       string  `json:"task_id"`
}

const (
	StatusProcessing  = "Processing..."
	StatusNotFound    = "Not found"
	StatusSuccess     = "Success"
	StatusError       = "Error scoring"
	StatusWrongTaskID = "Wrong task id"
)

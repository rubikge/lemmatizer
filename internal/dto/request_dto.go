package dto

type RequestData struct {
	Message string  `json:"message"`
	Product Product `json:"product"`
}

package models

type SearchProduct struct {
	ProductTitle     string
	RequiredKeywords []string
	MinCountWords    int
	SearchKeywords   []SearchKeyword
}

type SearchKeyword struct {
	Word              string
	Weight            float64
	RequiredWordIndex int
}

type SearchResult struct {
	Found        bool    `json:"found"`
	ProductTitle string  `json:"product_title"`
	TotalScore   float64 `json:"total_score"`
}

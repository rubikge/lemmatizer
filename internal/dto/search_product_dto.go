package dto

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

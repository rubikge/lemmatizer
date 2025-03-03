package dto

type Product struct {
	CommonPositiveKeywords []string     `json:"common_positive_keywords"`
	CommonNegativeKeywords []string     `json:"common_negative_keywords"`
	SubProducts            []SubProduct `json:"products"`
}

type SubProduct struct {
	ProductTitle         string                `json:"product_title"`
	RequiredKeywords     []string              `json:"required_keywords"`
	MinCountWords        int                   `json:"min_count_words"`
	Weight               float64               `json:"waight"`
	KeywordsWithSynonyms []KeywordWithSynonyms `json:"keywords_with_synonyms"`
}

type KeywordWithSynonyms struct {
	Synonyms []string `json:"sinonyms"`
	Weight   float64  `json:"waight"`
}

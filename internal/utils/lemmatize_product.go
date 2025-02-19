package utils

// import (
// 	"strings"

// 	"github.com/rubikge/lemmatizer/internal/models"
// 	"github.com/rubikge/lemmatizer/internal/services"
// )

// func LemmatizeProduct(product *models.Product, s *services.LemmatizerService) error {
// 	text := processProductJson(product)

// 	lemmas, err := s.GetLemmasArray(text)
// 	if err != nil {
// 		return err
// 	}

// 	updateProduct(product, lemmas)
// 	return nil
// }

// func processProductJson(product *models.Product) string {
// 	var sb strings.Builder

// 	sb.WriteString(strings.Join(product.CommonPositiveKeywords, " ") + " ")
// 	sb.WriteString(strings.Join(product.CommonNegativeKeywords, " ") + " ")

// 	for _, p := range product.SubProducts {
// 		sb.WriteString(strings.Join(p.RequiredKeywords, " ") + " ")
// 		for _, synonymGroup := range p.KeywordsWithSynonyms {
// 			sb.WriteString(strings.Join(synonymGroup.Synonyms, " ") + " ")
// 		}
// 	}

// 	return sb.String()
// }

// func updateProduct(product *models.Product, lemmas []services.Lemma) error {
// 	updatedProduct := *product
// 	lemmasIndex := 0

// 	filteredKeywords := updatedProduct.CommonPositiveKeywords[:0]
// 	seen := make(map[string]struct{})

// 	for _, keyword := range updatedProduct.CommonPositiveKeywords {
// 		if !(lemmasIndex < len(lemmas) && keyword == lemmas[lemmasIndex].Word) {
// 			return
// 		}
// 		{
// 			if lemmas[lemmasIndex].Lemma == "" {
// 				lemmasIndex++
// 				continue
// 			}
// 			keyword = lemmas[lemmasIndex].Lemma
// 			lemmasIndex++
// 		}

// 		if _, exists := seen[keyword]; !exists {
// 			seen[keyword] = struct{}{}
// 			filteredKeywords = append(filteredKeywords, keyword)
// 		}
// 	}

// 	updatedProduct.CommonPositiveKeywords = filteredKeywords
// }

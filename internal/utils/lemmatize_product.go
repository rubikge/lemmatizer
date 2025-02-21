package utils

import (
	"fmt"
	"strings"

	"github.com/rubikge/lemmatizer/internal/models"
	"github.com/rubikge/lemmatizer/internal/services"
)

func GetLemmatizedSearchProduct(product *models.Product, lemmatizer *services.LemmatizerService) ([]models.SearchProduct, error) {
	searchProducts, err := getRawSearchProducts(product)
	if err != nil {
		return nil, err
	}

	err = lemmatizeSearchProducts(&searchProducts, lemmatizer)
	if err != nil {
		return nil, err
	}

	return searchProducts, nil
}

func getRawSearchProducts(product *models.Product) ([]models.SearchProduct, error) {
	const positiveSearchKeywordWeight = 0.95
	var searchProducts []models.SearchProduct

	for _, subProduct := range product.SubProducts {
		var searchKeywords []models.SearchKeyword

		// Convert required keywords indexes into a set to check JSON
		requiredKeywordIndexSet := make(map[int]struct{}, len(subProduct.RequiredKeywords))
		for i := range subProduct.RequiredKeywords {
			requiredKeywordIndexSet[i] = struct{}{}
		}

		// Adding synonym groups
		for _, keyword := range subProduct.KeywordsWithSynonyms {
			var synonymKeywords []models.SearchKeyword

			// Create a map to avoid duplicate synonyms
			synonymSet := make(map[string]struct{}, len(keyword.Synonyms))
			for _, synonym := range keyword.Synonyms {
				synonymSet[synonym] = struct{}{}

				synonymKeywords = append(synonymKeywords, models.SearchKeyword{
					Word:              synonym,
					Weight:            keyword.Weight,
					RequiredWordIndex: -1,
				})
			}

			// Check if any synonym is a required keyword
			for requiredKeywordIndex := range requiredKeywordIndexSet {
				if _, isRequired := synonymSet[subProduct.RequiredKeywords[requiredKeywordIndex]]; isRequired {
					for i := range synonymKeywords {
						synonymKeywords[i].RequiredWordIndex = requiredKeywordIndex
					}
					delete(requiredKeywordIndexSet, requiredKeywordIndex)
					break
				}
			}

			searchKeywords = append(searchKeywords, synonymKeywords...)
		}

		if len(requiredKeywordIndexSet) != 0 {
			return nil, fmt.Errorf("JSON error, Requred keyword for '%s' missing in synonyms", subProduct.ProductTitle)
		}

		// Get positive keywords
		positiveSearchKeywords := make([]models.SearchKeyword, len(product.CommonPositiveKeywords))
		for i, keyword := range product.CommonPositiveKeywords {
			positiveSearchKeywords[i] = models.SearchKeyword{
				Word:              keyword,
				Weight:            positiveSearchKeywordWeight,
				RequiredWordIndex: -1,
			}
		}

		// add positive keywords at the end
		searchKeywords = append(searchKeywords, positiveSearchKeywords...)

		searchProducts = append(searchProducts, models.SearchProduct{
			ProductTitle:     subProduct.ProductTitle,
			RequiredKeywords: subProduct.RequiredKeywords,
			MinCountWords:    subProduct.MinCountWords,
			SearchKeywords:   searchKeywords,
		})
	}

	return searchProducts, nil
}

func lemmatizeSearchProducts(searchProducts *[]models.SearchProduct, lemmatizer *services.LemmatizerService) error {
	var textBuilder strings.Builder

	// Collect all words into a single text string
	for _, searchProduct := range *searchProducts {
		for _, searchKeyword := range searchProduct.SearchKeywords {
			word := searchKeyword.Word
			if strings.Contains(word, " ") {
				word = "errr"
			}
			textBuilder.WriteString(word)
			textBuilder.WriteString(" ")
		}
	}

	lemmas, err := lemmatizer.GetLemmasArray(textBuilder.String())
	if err != nil {
		return err
	}

	if len(lemmas) < len(*searchProducts) {
		return fmt.Errorf("lemmatization error: lemmas count is less than expected")
	}

	tempSearchProducts := make([]models.SearchProduct, len(*searchProducts))
	copy(tempSearchProducts, *searchProducts)

	lemmasIndex := 0

	for i := range tempSearchProducts {
		filteredSearchKeywords := make([]models.SearchKeyword, 0, len(tempSearchProducts[i].SearchKeywords))

		for j := range tempSearchProducts[i].SearchKeywords {
			if lemmasIndex >= len(lemmas) {
				return fmt.Errorf("lemmatization error: lemma index out of range")
			}

			lemma := lemmas[lemmasIndex].Lemma
			lemmasIndex++

			if lemma == "" || lemma == "errr" {
				continue
			}

			// Update the keyword
			tempSearchProducts[i].SearchKeywords[j].Word = lemma
			filteredSearchKeywords = append(filteredSearchKeywords, tempSearchProducts[i].SearchKeywords[j])
		}

		// Update slice with filtered keywords
		tempSearchProducts[i].SearchKeywords = filteredSearchKeywords
	}

	// Final lemma count check
	if lemmasIndex != len(lemmas) {
		return fmt.Errorf("lemmatization error: lemma count mismatch")
	}

	*searchProducts = tempSearchProducts

	return nil
}

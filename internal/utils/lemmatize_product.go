package utils

import (
	"fmt"
	"strings"

	"github.com/rubikge/lemmatizer/internal/models"
	"github.com/rubikge/lemmatizer/internal/services"
)

func GetLemmatizedSearchProduct(product *models.Product, lemmatizer *services.LemmatizerService) ([]models.SearchProduct, error) {
	searchProducts := getRawSearchProducts(product)

	err := lemmatizeSearchProducts(&searchProducts, lemmatizer)
	if err != nil {
		return nil, err
	}

	return searchProducts, nil
}

func getRawSearchProducts(product *models.Product) []models.SearchProduct {
	const positiveSearchKeywordWeight = 0.95
	var searchProducts []models.SearchProduct

	// Getting positive keywords
	positiveSearchKeywords := make([]models.SearchKeyword, len(product.CommonPositiveKeywords))
	for i, keyword := range product.CommonPositiveKeywords {
		positiveSearchKeywords[i] = models.SearchKeyword{
			Word:   keyword,
			Weight: positiveSearchKeywordWeight,
		}
	}

	for _, subProduct := range product.SubProducts {
		// Create a copy of positiveSearchKeywords
		searchKeywords := make([]models.SearchKeyword, len(positiveSearchKeywords))
		copy(searchKeywords, positiveSearchKeywords)

		// Convert required keywords into a set for faster lookup
		requiredKeywordSet := make(map[string]struct{}, len(subProduct.RequiredKeywords))
		for _, requiredKeyword := range subProduct.RequiredKeywords {
			requiredKeywordSet[requiredKeyword] = struct{}{}
		}

		// Adding synonym groups
		for _, keyword := range subProduct.KeywordsWithSynonyms {
			var synonymKeywords []models.SearchKeyword

			// Create a map to avoid duplicate synonyms
			synonymSet := make(map[string]struct{}, len(keyword.Synonyms))
			for _, synonym := range keyword.Synonyms {
				synonymSet[synonym] = struct{}{}

				synonymKeywords = append(synonymKeywords, models.SearchKeyword{
					Word:   synonym,
					Weight: keyword.Weight,
				})
			}

			// Check if any synonym is a required keyword
			for _, requiredKeyword := range subProduct.RequiredKeywords {
				if _, isRequired := synonymSet[requiredKeyword]; isRequired {
					for i := range synonymKeywords {
						synonymKeywords[i].RequiredWord = requiredKeyword
					}
					break
				}
			}

			searchKeywords = append(searchKeywords, synonymKeywords...)
		}

		searchProducts = append(searchProducts, models.SearchProduct{
			ProductTitle:           subProduct.ProductTitle,
			RequiredKeywordsNumber: len(subProduct.RequiredKeywords),
			MinCountWords:          subProduct.MinCountWords,
			SearchKeywords:         searchKeywords,
		})
	}

	return searchProducts
}

func lemmatizeSearchProducts(searchProducts *[]models.SearchProduct, lemmatizer *services.LemmatizerService) error {
	var textBuilder strings.Builder

	// Collect all words into a single text string
	for _, searchProduct := range *searchProducts {
		for _, searchKeyword := range searchProduct.SearchKeywords {
			textBuilder.WriteString(searchKeyword.Word)
			textBuilder.WriteString(" ")
		}
	}

	lemmas, err := lemmatizer.GetLemmasArray(textBuilder.String())
	if err != nil {
		return err
	}

	if len(lemmas) < len(*searchProducts) {
		return fmt.Errorf("Lemmatization error: Lemmas count is less than expected")
	}

	tempSearchProducts := make([]models.SearchProduct, len(*searchProducts))
	copy(tempSearchProducts, *searchProducts)

	lemmasIndex := 0

	for i := range tempSearchProducts {
		filteredSearchKeywords := make([]models.SearchKeyword, 0, len(tempSearchProducts[i].SearchKeywords))

		for j := range tempSearchProducts[i].SearchKeywords {
			if lemmasIndex >= len(lemmas) {
				return fmt.Errorf("Lemmatization error: Lemma index out of range")
			}

			lemma := lemmas[lemmasIndex].Lemma
			lemmasIndex++

			if lemma == "" {
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
		return fmt.Errorf("Lemmatization error: Lemma count mismatch")
	}

	*searchProducts = tempSearchProducts

	return nil
}

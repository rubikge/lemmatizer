package lemmatizer

import (
	"fmt"
	"slices"
	"strings"

	"github.com/rubikge/lemmatizer/internal/dto"
	"github.com/rubikge/lemmatizer/internal/mystem"
)

type LemmatizerService struct {
	repo *mystem.MystemRepository
}

func NewLemmatizerService(repo *mystem.MystemRepository) *LemmatizerService {
	return &LemmatizerService{repo: repo}
}

func (ls *LemmatizerService) GetLemmasArray(text string) ([]Lemma, error) {
	analysis, err := ls.repo.GetAnalysis(text)
	if err != nil {
		return nil, err
	}

	var lemmasArray []Lemma

	for _, word := range analysis {
		if len(word.Analysis) == 0 {
			lemmasArray = append(lemmasArray, Lemma{Word: word.Text, Lemma: word.Text})
			continue
		}

		analysis := word.Analysis[0]

		lemma := ""
		if slices.ContainsFunc(mystem.NeededPrefixes, func(prefix string) bool {
			return strings.HasPrefix(analysis.Gr, prefix+"=") || strings.HasPrefix(analysis.Gr, prefix+",")
		}) {
			lemma = analysis.Lex
		}
		lemmasArray = append(lemmasArray, Lemma{Word: word.Text, Lemma: lemma})
	}

	return lemmasArray, nil
}

func (ls *LemmatizerService) GetLemmas(text string) ([]string, error) {
	fmt.Printf("Lemmatizing string '%s'...\n", text)
	lemmas, err := ls.GetLemmasArray(text)
	if err != nil {
		return nil, err
	}

	words := make([]string, 0, len(lemmas))
	for _, lemma := range lemmas {
		word := lemma.Lemma
		if word != "" {
			words = append(words, lemma.Lemma)
		}
	}
	fmt.Printf("Lemmas: %s.\n\n", strings.Join(words, ", "))
	return words, nil
}

func (ls *LemmatizerService) GetLemmatizedSearchProduct(product *dto.Product) ([]dto.SearchProduct, error) {
	searchProducts, err := getRawSearchProducts(product, positiveSearchKeywordWeight)
	if err != nil {
		return nil, err
	}

	err = ls.lemmatizeSearchProducts(&searchProducts)
	if err != nil {
		return nil, err
	}

	return searchProducts, nil
}

func getRawSearchProducts(product *dto.Product, positiveWeight float64) ([]dto.SearchProduct, error) {
	var searchProducts []dto.SearchProduct

	for _, subProduct := range product.SubProducts {
		var searchKeywords []dto.SearchKeyword

		// Convert required keywords indexes into a set to check JSON
		requiredKeywordIndexSet := make(map[int]struct{}, len(subProduct.RequiredKeywords))
		for i := range subProduct.RequiredKeywords {
			requiredKeywordIndexSet[i] = struct{}{}
		}

		// Adding synonym groups
		for _, keyword := range subProduct.KeywordsWithSynonyms {
			var synonymKeywords []dto.SearchKeyword

			// Create a map to avoid duplicate synonyms
			synonymSet := make(map[string]struct{}, len(keyword.Synonyms))
			for _, synonym := range keyword.Synonyms {
				synonymSet[synonym] = struct{}{}

				synonymKeywords = append(synonymKeywords, dto.SearchKeyword{
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
		positiveSearchKeywords := make([]dto.SearchKeyword, len(product.CommonPositiveKeywords))
		for i, keyword := range product.CommonPositiveKeywords {
			positiveSearchKeywords[i] = dto.SearchKeyword{
				Word:              keyword,
				Weight:            positiveWeight,
				RequiredWordIndex: -1,
			}
		}

		// add positive keywords at the end
		searchKeywords = append(searchKeywords, positiveSearchKeywords...)

		searchProducts = append(searchProducts, dto.SearchProduct{
			ProductTitle:     subProduct.ProductTitle,
			RequiredKeywords: subProduct.RequiredKeywords,
			MinCountWords:    subProduct.MinCountWords,
			SearchKeywords:   searchKeywords,
		})
	}

	return searchProducts, nil
}

func (ls *LemmatizerService) lemmatizeSearchProducts(searchProducts *[]dto.SearchProduct) error {
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

	lemmas, err := ls.GetLemmasArray(textBuilder.String())
	if err != nil {
		return err
	}

	if len(lemmas) < len(*searchProducts) {
		return fmt.Errorf("lemmatization error: lemmas count is less than expected")
	}

	tempSearchProducts := make([]dto.SearchProduct, len(*searchProducts))
	copy(tempSearchProducts, *searchProducts)

	lemmasIndex := 0

	for i := range tempSearchProducts {
		filteredSearchKeywords := make([]dto.SearchKeyword, 0, len(tempSearchProducts[i].SearchKeywords))

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

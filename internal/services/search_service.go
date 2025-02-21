package services

import (
	"fmt"
	"log"
	"strconv"

	"github.com/blevesearch/bleve"
	"github.com/rubikge/lemmatizer/internal/models"
)

type searchKeyword struct {
	Word string `json:"word"`
}

func getIndex(searchKeywords *[]models.SearchKeyword) (bleve.Index, error) {
	indexMapping := bleve.NewIndexMapping()
	docMapping := bleve.NewDocumentMapping()

	wordFieldMapping := bleve.NewTextFieldMapping()
	wordFieldMapping.Store = true
	wordFieldMapping.Index = true
	docMapping.AddFieldMappingsAt("word", wordFieldMapping)

	indexMapping.AddDocumentMapping("search_keyword", docMapping)

	index, err := bleve.NewMemOnly(indexMapping)
	if err != nil {
		return nil, err
	}

	for i, searchKeyWord := range *searchKeywords {
		doc := searchKeyword{Word: searchKeyWord.Word}
		docID := strconv.Itoa(i)
		err = index.Index(docID, doc)
		if err != nil {
			return nil, err
		}
	}

	return index, nil
}

func GetScore(words *[]string, searchProducts *[]models.SearchProduct) models.SearchResult {
	const goalTotalScore = 1
	const minRequiredKeywordsCount = 1

	for _, searchProduct := range *searchProducts {
		if len(*words) < searchProduct.MinCountWords {
			continue
		}

		index, err := getIndex(&searchProduct.SearchKeywords)
		if err != nil {
			log.Println("Error getting index:", err)
			continue
		}

		totalScore := 0.0

		requiredKeywordIndexSet := make(map[int]struct{}, len(searchProduct.RequiredKeywords))
		for i := range searchProduct.RequiredKeywords {
			requiredKeywordIndexSet[i] = struct{}{}
		}

		for _, word := range *words {
			query := bleve.NewFuzzyQuery(word)
			query.Fuzziness = 2
			query.SetField("word")
			searchReq := bleve.NewSearchRequest(query)
			searchReq.Size = 1

			searchRes, err := index.Search(searchReq)
			if err != nil {
				log.Println("Search error: ", err)
				continue
			}

			if len(searchRes.Hits) == 0 {
				continue
			}

			keywordIndex, err := strconv.Atoi(searchRes.Hits[0].ID)
			if err != nil || keywordIndex >= len(searchProduct.SearchKeywords) {
				log.Println("Index error: ", err)
				continue
			}

			keyword := searchProduct.SearchKeywords[keywordIndex]

			totalScore += keyword.Weight
			if keyword.RequiredWordIndex != -1 {
				delete(requiredKeywordIndexSet, keyword.RequiredWordIndex)
			}

			fmt.Println(models.SearchResult{
				ProductTitle: searchProduct.ProductTitle,
				TotalScore:   totalScore,
			})

			if totalScore >= goalTotalScore && len(searchProduct.RequiredKeywords)-len(requiredKeywordIndexSet) <= minRequiredKeywordsCount {
				return models.SearchResult{
					ProductTitle: searchProduct.ProductTitle,
					TotalScore:   totalScore,
				}
			}
		}
	}

	return models.SearchResult{}
}

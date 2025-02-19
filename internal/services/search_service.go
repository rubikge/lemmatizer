package services

import (
	"fmt"
	"log"

	"github.com/blevesearch/bleve"
	"github.com/rubikge/lemmatizer/internal/models"
)

func GetTotalScore(words []string, searchData *models.SearchData) float64 {
	index, err := createIndex(searchData)
	if err != nil {
		log.Fatal("Error creating index:", err)
	}

	totalScore := 0.0

	for _, word := range words {
		query := bleve.NewFuzzyQuery(word)
		query.Fuzziness = 2
		searchReq := bleve.NewSearchRequest(query)
		searchReq.Fields = []string{"word", "weight"}
		searchRes, err := index.Search(searchReq)
		if err != nil {
			log.Println("Search error: ", err)
			continue
		}

		for _, hit := range searchRes.Hits {
			fmt.Println(hit)
			if weight, ok := hit.Fields["weight"].(float64); ok {
				totalScore += weight
			}
		}
	}
	defer index.Close()
	return totalScore
}

func createIndex(searchData *models.SearchData) (bleve.Index, error) {
	mapping := bleve.NewIndexMapping()
	index, err := bleve.NewMemOnly(mapping)
	if err != nil {
		return nil, err
	}

	for _, keyword := range searchData.Keywords {
		document := map[string]interface{}{
			"word":   keyword.Word,
			"weight": keyword.Weight,
		}
		err := index.Index(keyword.Word, document)
		if err != nil {
			return nil, err
		}
	}
	return index, nil
}

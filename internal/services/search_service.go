package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/blevesearch/bleve"
	"github.com/rubikge/lemmatizer/internal/models"
)

type searchKeyword struct {
	Word string `json:"word"`
}

const minRequiredWordsCount = 2
const goalTotalScore = 1.0

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
	fmt.Printf(
		"Goal Total Score - %.2f, Min required words number - %d\n\n",
		goalTotalScore,
		minRequiredWordsCount,
	)
	for _, searchProduct := range *searchProducts {
		fmt.Printf(
			"Calculating the score of subproduct %s...\nRequired words: %s.\n",
			searchProduct.ProductTitle,
			strings.Join(searchProduct.RequiredKeywords, ", "),
		)

		if wordsNum := len(*words); wordsNum < searchProduct.MinCountWords {
			fmt.Printf("Message length (%d) is less than min_count_words(%d)", wordsNum, searchProduct.MinCountWords)
			continue
		}

		index, err := getIndex(&searchProduct.SearchKeywords)
		if err != nil {
			fmt.Println("Error getting index:", err)
			continue
		}

		totalScore := 0.0
		requiredWordsCount := 0

		requiredKeywordIndexSet := make(map[int]struct{}, len(searchProduct.RequiredKeywords))
		for i := range searchProduct.RequiredKeywords {
			requiredKeywordIndexSet[i] = struct{}{}
		}

		fmt.Printf("Match words:\n")

		for _, word := range *words {
			query := bleve.NewFuzzyQuery(word)
			query.Fuzziness = 2
			query.SetField("word")
			searchReq := bleve.NewSearchRequest(query)
			searchReq.Size = 1

			searchRes, err := index.Search(searchReq)
			if err != nil {
				fmt.Println("Search error: ", err)
				continue
			}

			if len(searchRes.Hits) == 0 {
				continue
			}

			keywordIndex, err := strconv.Atoi(searchRes.Hits[0].ID)
			if err != nil || keywordIndex >= len(searchProduct.SearchKeywords) {
				fmt.Println("Index error: ", err)
				continue
			}

			keyword := searchProduct.SearchKeywords[keywordIndex]

			totalScore += keyword.Weight
			fmt.Printf("     %s -> { %s, %.2f", word, keyword.Word, keyword.Weight)

			if keyword.RequiredWordIndex != -1 {
				delete(requiredKeywordIndexSet, keyword.RequiredWordIndex)
				requiredWordsCount++
				fmt.Print(", required")
			}
			fmt.Println(" }")

			if totalScore >= goalTotalScore && requiredWordsCount >= minRequiredWordsCount {
				fmt.Printf("Result: positive. Total score - %.2f.\n\n\n", totalScore)
				return models.SearchResult{
					ProductTitle: searchProduct.ProductTitle,
					TotalScore:   totalScore,
				}
			}
		}
		fmt.Printf("Result: negative. Total score - %.2f, required words number - %d.\n\n\n", totalScore, requiredWordsCount)
	}

	return models.SearchResult{}
}

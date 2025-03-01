package services

import (
	"fmt"
	"slices"
	"strings"

	"github.com/rubikge/lemmatizer/internal/models"
	"github.com/rubikge/lemmatizer/internal/mystem"
	"github.com/rubikge/lemmatizer/internal/repository"
)

type LemmatizerService struct {
	repo *repository.MystemRepository
}

func NewLemmatizerService(repo *repository.MystemRepository) *LemmatizerService {
	return &LemmatizerService{repo: repo}
}

func (s *LemmatizerService) GetLemmasArray(text string) ([]models.Lemma, error) {
	analysis, err := s.repo.GetAnalysis(text)
	if err != nil {
		return nil, err
	}

	var lemmasArray []models.Lemma

	for _, word := range analysis {
		if len(word.Analysis) == 0 {
			lemmasArray = append(lemmasArray, models.Lemma{Word: word.Text, Lemma: word.Text})
			continue
		}

		analysis := word.Analysis[0]

		lemma := ""
		if slices.ContainsFunc(mystem.NeededPrefixes, func(prefix string) bool {
			return strings.HasPrefix(analysis.Gr, prefix+"=") || strings.HasPrefix(analysis.Gr, prefix+",")
		}) {
			lemma = analysis.Lex
		}
		lemmasArray = append(lemmasArray, models.Lemma{Word: word.Text, Lemma: lemma})
	}

	return lemmasArray, nil
}

func (s *LemmatizerService) GetLemmas(text string) ([]string, error) {
	fmt.Printf("Lemmatizing string '%s'...\n", text)
	lemmas, err := s.GetLemmasArray(text)
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

// func (s *LemmatizerService) GetLemmasMap(text string) (map[string]struct{}, error) {
// 	lemmasArray, err := s.GetLemmasArray(text)
// 	if err != nil {
// 		return nil, err
// 	}

// 	lemmasMap := make(map[string]struct{})

// 	for _, lemma := range lemmasArray {
// 		lemmasMap[lemma] = struct{}{}
// 	}

// 	return lemmasMap, nil
// }

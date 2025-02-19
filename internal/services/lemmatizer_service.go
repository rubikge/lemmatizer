package services

import (
	"slices"
	"strings"

	"github.com/rubikge/lemmatizer/internal/mystem"
	"github.com/rubikge/lemmatizer/internal/repository"
)

type LemmatizerService struct {
	repo *repository.MystemRepository
}

func NewLemmatizerService(repo *repository.MystemRepository) *LemmatizerService {
	return &LemmatizerService{repo: repo}
}

type Lemma struct {
	Word, Lemma string
}

func (s *LemmatizerService) GetLemmasArray(text string) ([]Lemma, error) {
	wordStream, err := s.repo.GetDataStream(text)
	if err != nil {
		return nil, err
	}

	var lemmasArray []Lemma

	for word := range wordStream {
		if len(word.Analysis) == 0 {
			lemmasArray = append(lemmasArray, Lemma{word.Text, word.Text})
			continue
		}

		analysis := word.Analysis[0]

		lemma := ""
		if slices.ContainsFunc(mystem.NeededPrefixes, func(prefix string) bool {
			return strings.HasPrefix(analysis.Gr, prefix+"=") || strings.HasPrefix(analysis.Gr, prefix+",")
		}) {
			lemma = analysis.Lex
		}
		lemmasArray = append(lemmasArray, Lemma{word.Text, lemma})
	}

	return lemmasArray, nil
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

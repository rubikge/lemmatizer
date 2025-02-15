package service

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

func (s *LemmatizerService) ProcessData(text string) ([]string, error) {
	wordStream, err := s.repo.GetDataStream()
	if err != nil {
		return nil, err
	}

	lemmasSet := make(map[string]struct{})

	for word := range wordStream {
		if len(word.Analysis) == 0 {
			lemmasSet[word.Text] = struct{}{}
			continue
		}

		analysis := word.Analysis[0]

		if slices.ContainsFunc(mystem.NeededPrefixes, func(prefix string) bool {
			return strings.HasPrefix(analysis.Gr, prefix+"=") || strings.HasPrefix(analysis.Gr, prefix+",")
		}) {
			lemmasSet[analysis.Lex] = struct{}{}
		}
	}

	var lemmas []string
	for item := range lemmasSet {
		lemmas = append(lemmas, item)
	}

	return lemmas, nil
}

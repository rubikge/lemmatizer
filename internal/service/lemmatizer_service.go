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

	var lemmas []string

	for word := range wordStream {
		if len(word.Analysis) == 0 {
			lemmas = append(lemmas, word.Text)
			continue
		}

		analysis := word.Analysis[0]

		if slices.ContainsFunc(mystem.NeededPrefixes, func(prefix string) bool {
			return strings.HasPrefix(analysis.Gr, prefix+"=") || strings.HasPrefix(analysis.Gr, prefix+",")
		}) {
			lemmas = append(lemmas, analysis.Lex)
		}
	}

	return lemmas, nil
}

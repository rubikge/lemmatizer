package repository

import (
	"encoding/json"
	"fmt"

	"github.com/rubikge/lemmatizer/internal/model"
)

type MystemRepository struct {
	data []byte
}

func NewMystemRepository(data []byte) *MystemRepository {
	return &MystemRepository{data: data}
}

func (r *MystemRepository) GetData() ([]model.AnalizedWord, error) {
	var words []model.AnalizedWord
	err := json.Unmarshal(r.data, &words)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return lemmas, nil
}

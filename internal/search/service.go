package search

import (
	"fmt"

	"github.com/rubikge/lemmatizer/internal/dto"
	"github.com/rubikge/lemmatizer/internal/lemmatizer"
	"github.com/rubikge/lemmatizer/internal/positive_check"
	"github.com/rubikge/lemmatizer/internal/query_scorer"
	"github.com/rubikge/lemmatizer/internal/redis"
)

type Service struct {
	name       string
	positive   *positive_check.Service
	Lemmatizer *lemmatizer.LemmatizerService
	redisQueue *redis.RedisQueue
}

func NewService(name string) (*Service, error) {
	positive, err := positive_check.NewService()
	if err != nil {
		return nil, err
	}

	qs := query_scorer.NewService()

	rs, err := redis.NewRedisQueue(qs)
	if err != nil {
		return nil, err
	}

	service := &Service{
		name:       name,
		positive:   positive,
		Lemmatizer: qs.Lemmatizer,
		redisQueue: rs,
	}

	// Start Redis worker with the provided name
	if err := rs.StartWorker(name); err != nil {
		return nil, fmt.Errorf("failed to start worker with name %s: %w", name, err)
	}

	return service, nil
}

// 2-step ProcessSearch processes the search request and returns the initial response
func (s *Service) ProcessSearch(request *dto.RequestData) (*dto.SearchResult, error) {
	// 1st step. Do positive check using query scorer
	containsPositiveWords := s.positive.ContainsPositiveWords(request.Message)

	if !containsPositiveWords {
		return &dto.SearchResult{
			Status: dto.StatusNotFound,
		}, nil
	}

	// 2nd step. Send to Redis queue for processing
	taskID, err := s.redisQueue.AddRequestToQueue(request)
	if err != nil {
		return nil, err
	}

	// Return Processing status
	return &dto.SearchResult{
		Status: dto.StatusProcessing,
		TaskID: taskID,
	}, nil
}

// GetResult retrieves the search result from Redis
func (s *Service) GetResult(taskID string) (*dto.SearchResult, error) {
	result, err := s.redisQueue.GetResponseFromQueue(taskID)

	return result, err
}

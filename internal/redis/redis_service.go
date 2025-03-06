package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rubikge/lemmatizer/internal/dto"
	"github.com/rubikge/lemmatizer/internal/query_scorer"
)

type RedisQueue struct {
	rdb         *redis.Client
	ctx         context.Context
	workers     map[string]struct{}
	queryScorer *query_scorer.Service
}

func NewRedisQueue(qs *query_scorer.Service) (*RedisQueue, error) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: Addr,
	})

	err := rdb.XGroupCreateMkStream(ctx, StreamName, ConsumerGroup, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &RedisQueue{
		rdb:         rdb,
		ctx:         ctx,
		workers:     map[string]struct{}{},
		queryScorer: qs,
	}, nil
}

func (rq *RedisQueue) AddRequestToQueue(requestData *dto.RequestData) (string, error) {
	data, err := json.Marshal(requestData)
	if err != nil {
		return "", err
	}

	taskID, err := rq.rdb.XAdd(rq.ctx, &redis.XAddArgs{
		Stream: StreamName,
		Values: map[string]interface{}{
			"data": data,
		},
	}).Result()

	if err != nil {
		return "", err
	}
	return taskID, nil
}

func (rq *RedisQueue) GetResponseFromQueue(taskID string) (*dto.SearchResult, error) {
	result, err := rq.rdb.Get(rq.ctx, taskID).Result()

	if err == redis.Nil {
		// Check if task exists in stream
		messages, err1 := rq.rdb.XRange(rq.ctx, StreamName, taskID, taskID).Result()
		if err1 != nil {
			return nil, err
		}

		if len(messages) == 0 {
			return &dto.SearchResult{
				Status: dto.StatusWrongTaskID,
				TaskID: taskID,
			}, nil
		}

		return &dto.SearchResult{
			Status: dto.StatusProcessing,
			TaskID: taskID,
		}, nil
	}

	if err != nil {
		return nil, err
	}

	var searchResult dto.SearchResult
	err = json.Unmarshal([]byte(result), &searchResult)
	if err != nil {
		return nil, err
	}
	searchResult.TaskID = taskID

	return &searchResult, nil
}

func (rq *RedisQueue) StartWorker(workerName string) error {
	if _, isStarted := rq.workers[workerName]; isStarted {
		return fmt.Errorf("worker %s is already started", workerName)
	}

	rq.workers[workerName] = struct{}{}

	go func() {
		defer delete(rq.workers, workerName)
		rq.startNewWorker(workerName)
	}()

	return nil
}

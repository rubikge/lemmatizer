package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// RedisQueueInterface defines the interface for queue operations
type RedisQueueInterface interface {
	AddRequestToQueue(data string) (string, error)
	GetResponseFromQueue(taskId string) (string, error)
	StartWorker(workerName string, process func(string) (string, error)) error
}

type RedisQueue struct {
	rdb     *redis.Client
	ctx     context.Context
	workers map[string]struct{}
}

func NewRedisQueue() (*RedisQueue, error) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	err := rdb.XGroupCreateMkStream(ctx, streamName, consumerGroup, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &RedisQueue{
		rdb:     rdb,
		ctx:     ctx,
		workers: map[string]struct{}{},
	}, nil
}

func (rq *RedisQueue) AddRequestToQueue(data string) (string, error) {
	taskId, err := rq.rdb.XAdd(rq.ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: map[string]interface{}{
			"data": data,
		},
	}).Result()

	if err != nil {
		return "", err
	}
	return taskId, nil
}

func (rq *RedisQueue) GetResponseFromQueue(taskId string) (string, error) {
	// First check if we have a result
	result, err := rq.rdb.Get(rq.ctx, taskId).Result()

	if err == redis.Nil {
		// Check if task exists in stream
		messages, err1 := rq.rdb.XRange(rq.ctx, streamName, taskId, taskId).Result()
		if err1 != nil || len(messages) == 0 {
			response := map[string]interface{}{
				"status": StatusError,
				"error":  "message not found",
			}
			resultJSON, _ := json.Marshal(response)
			return string(resultJSON), nil
		}

		// Task exists but still processing
		response := map[string]interface{}{
			"status": StatusProcessing,
		}
		resultJSON, _ := json.Marshal(response)
		return string(resultJSON), nil
	}

	if err != nil {
		response := map[string]interface{}{
			"status": StatusError,
			"error":  err.Error(),
		}
		resultJSON, _ := json.Marshal(response)
		return string(resultJSON), nil
	}

	if result == "" {
		response := map[string]interface{}{
			"status": StatusError,
			"error":  "empty result",
		}
		resultJSON, _ := json.Marshal(response)
		return string(resultJSON), nil
	}

	// Result already contains proper JSON structure from worker
	return result, nil
}

func (rq *RedisQueue) StartWorker(workerName string, process func(string) (string, error)) error {
	if _, isStarted := rq.workers[workerName]; isStarted {
		return fmt.Errorf("worker %s is already started", workerName)
	}

	rq.workers[workerName] = struct{}{}

	go func() {
		defer delete(rq.workers, workerName)
		startNewWorker(rq, workerName, process)
	}()

	return nil
}

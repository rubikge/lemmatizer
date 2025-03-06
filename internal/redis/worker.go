package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rubikge/lemmatizer/internal/dto"
)

const (
	maxRetries = 3
	retryDelay = 5 * time.Second
)

func (rq *RedisQueue) startNewWorker(consumer string) {
	for {
		msgs, err := rq.rdb.XReadGroup(rq.ctx, &redis.XReadGroupArgs{
			Group:    ConsumerGroup,
			Consumer: consumer,
			Streams:  []string{StreamName, ">"},
			Count:    1,
			Block:    5 * time.Second,
		}).Result()

		if err == redis.Nil {
			continue
		} else if err != nil {
			fmt.Printf("Redis read error: %v\n", err)
			time.Sleep(retryDelay)
			continue
		}

		for _, stream := range msgs {
			for _, msg := range stream.Messages {
				taskID := msg.ID
				data := msg.Values["data"].(string)

				fmt.Printf("Processing task: %s\n", taskID)

				// Check if task has previous errors
				var taskError TaskError
				errorKey := fmt.Sprintf("error:%s", taskID)
				errorData, err := rq.rdb.Get(rq.ctx, errorKey).Result()
				if err != nil && err != redis.Nil {
					taskError = TaskError{
						Message:   "",
						Timestamp: time.Now().UTC().Format(time.RFC3339),
						Retries:   0,
					}
				} else if err == nil {
					if err := json.Unmarshal([]byte(errorData), &taskError); err != nil {
						fmt.Printf("Error unmarshaling task error: %v\n", err)
					}
				}

				// Process the task
				var requestData dto.RequestData
				if err := json.Unmarshal([]byte(data), &requestData); err != nil {
					fmt.Printf("Error unmarshaling data: %v\n", err)
					continue
				}
				result, err := rq.queryScorer.GetScore(&requestData)

				if err != nil {
					taskError.Retries++
					taskError.Message = err.Error()
					taskError.Timestamp = time.Now().UTC().Format(time.RFC3339)

					errorJSON, _ := json.Marshal(taskError)
					rq.rdb.Set(rq.ctx, errorKey, string(errorJSON), 24*time.Hour)

					if taskError.Retries >= maxRetries {
						fmt.Printf("Task %s failed after %d retries: %v\n", taskID, maxRetries, err)

						// Create error result
						result = &dto.SearchResult{
							Status: dto.StatusError,
						}

					} else {
						// Retry later - don't acknowledge the message
						fmt.Printf("Task %s failed, will retry (%d/%d): %v\n", taskID, taskError.Retries, maxRetries, err)
						time.Sleep(retryDelay)
						continue
					}
				}

				resultJSON, err := json.Marshal(*result)
				if err != nil {
					fmt.Printf("Error marshaling result: %v\n", err)
					continue
				}

				// Store the result
				rq.rdb.Set(rq.ctx, taskID, resultJSON, 24*time.Hour)
				rq.rdb.Del(rq.ctx, errorKey) // Clean up error tracking
				rq.rdb.XAck(rq.ctx, StreamName, ConsumerGroup, taskID)
			}
		}
	}
}

package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	maxRetries = 3
	retryDelay = 5 * time.Second
)

func startNewWorker(rq *RedisQueue, consumer string, process func(string) (string, error)) {
	for {
		msgs, err := rq.rdb.XReadGroup(rq.ctx, &redis.XReadGroupArgs{
			Group:    consumerGroup,
			Consumer: consumer,
			Streams:  []string{streamName, ">"},
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
				taskId := msg.ID
				data := msg.Values["data"].(string)

				fmt.Printf("Processing task: %s\n", taskId)

				// Check if task has previous errors
				var taskError TaskError
				errorKey := fmt.Sprintf("error:%s", taskId)
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
				result, err := process(data)

				if err != nil {
					taskError.Retries++
					taskError.Message = err.Error()
					taskError.Timestamp = time.Now().UTC().Format(time.RFC3339)

					errorJSON, _ := json.Marshal(taskError)
					rq.rdb.Set(rq.ctx, errorKey, string(errorJSON), 24*time.Hour)

					if taskError.Retries >= maxRetries {
						fmt.Printf("Task %s failed after %d retries: %v\n", taskId, maxRetries, err)

						// Create error response as raw JSON
						errorResponse := map[string]interface{}{
							"status": StatusError,
							"error":  err.Error(),
						}
						resultJSON, _ := json.Marshal(errorResponse)

						// Store final error state
						rq.rdb.Set(rq.ctx, taskId, string(resultJSON), 24*time.Hour)
						rq.rdb.XAck(rq.ctx, streamName, consumerGroup, taskId)
					} else {
						// Retry later - don't acknowledge the message
						fmt.Printf("Task %s failed, will retry (%d/%d): %v\n", taskId, taskError.Retries, maxRetries, err)
						time.Sleep(retryDelay)
						continue
					}
				} else {
					// Create response structure
					response := map[string]interface{}{
						"status": StatusSuccess,
					}

					// Check if result is valid JSON
					var resultData interface{}
					if err := json.Unmarshal([]byte(result), &resultData); err != nil {
						// If not valid JSON, use as raw string
						response["data"] = result
					} else {
						// If valid JSON, use parsed data
						response["data"] = resultData
					}

					// Marshal the final response
					wrappedResult, _ := json.Marshal(response)

					// Store the result
					rq.rdb.Set(rq.ctx, taskId, string(wrappedResult), 24*time.Hour)
					rq.rdb.Del(rq.ctx, errorKey) // Clean up error tracking
					rq.rdb.XAck(rq.ctx, streamName, consumerGroup, taskId)
				}
			}
		}
	}
}

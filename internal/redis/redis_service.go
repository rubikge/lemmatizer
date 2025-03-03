package redis

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	rdb     *redis.Client
	ctx     context.Context
	workers map[string]struct{}
}

func NewRedisQueue() *RedisQueue {
	ctx := context.Background()
	return &RedisQueue{
		rdb: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
		ctx:     ctx,
		workers: map[string]struct{}{},
	}
}

func (rq *RedisQueue) AddRequestToQueue(data string) (string, error) {
	taskId := strconv.Itoa(rand.Int())
	fmt.Println(taskId)

	err := rq.rdb.XAdd(rq.ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: map[string]interface{}{
			"data":   data,
			"taskId": taskId,
		},
	}).Err()

	if err != nil {
		return "", err
	}
	return taskId, nil
}

func (rq *RedisQueue) GetResponseFromQueue(taskId string) (*Response, error) {
	result, err := rq.rdb.Get(rq.ctx, taskId).Result()

	if err == redis.Nil {
		return &Response{
			Status: StatusProcessing,
		}, nil
	}

	if err != nil {
		return nil, err
	}

	if result == "" {
		return &Response{
			Status: StatusError,
			Data:   "",
		}, nil
	}

	return &Response{
		Status: StatusDone,
		Data:   result,
	}, nil
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

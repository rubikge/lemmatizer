package redis

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
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
			fmt.Println("Redis read error:", err)
			continue
		}

		for _, stream := range msgs {
			for _, msg := range stream.Messages {
				taskID := msg.ID
				data := msg.Values["data"].(string)

				fmt.Println("Processing task:", taskID)

				//update errors
				result, err := process(data)
				if err != nil {
					result = ""
				}

				rq.rdb.Set(rq.ctx, taskID, result, 10*time.Minute)

				rq.rdb.XAck(rq.ctx, streamName, consumerGroup, msg.ID)
			}
		}
	}
}

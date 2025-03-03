package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rubikge/lemmatizer/internal/redis"
)

type SearchController struct {
	rq *redis.RedisQueue
}

func NewSearchController(rq *redis.RedisQueue) *SearchController {
	return &SearchController{rq: rq}
}

func (c *SearchController) ProcessHandler(ctx fiber.Ctx) error {
	taskId, err := c.rq.AddRequestToQueue(string(ctx.Body()))
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": "Failed to enqueue task"})
	}

	return ctx.JSON(fiber.Map{"task_id": taskId})
}

func (c *SearchController) GetResultHandler(ctx fiber.Ctx) error {
	taskID := ctx.Params("taskID")
	result, err := c.rq.GetResponseFromQueue(taskID)

	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": "Redis error"})
	}

	return ctx.JSON(result)
}

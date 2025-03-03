package controller

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/rubikge/lemmatizer/internal/redis"
)

type SearchController struct {
	rq redis.RedisQueueInterface
}

func NewSearchController(rq redis.RedisQueueInterface) *SearchController {
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
	resultJSON, err := c.rq.GetResponseFromQueue(taskID)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"status": "Error",
			"error":  "Internal server error",
		})
	}

	// Parse the JSON response
	var response map[string]interface{}
	if err := json.Unmarshal([]byte(resultJSON), &response); err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"status": "Error",
			"error":  "Invalid response format",
		})
	}

	status, ok := response["status"].(string)
	if !ok {
		return ctx.Status(500).JSON(fiber.Map{
			"status": "Error",
			"error":  "Invalid response format",
		})
	}

	switch status {
	case redis.StatusError:
		return ctx.Status(500).JSON(response)
	case redis.StatusProcessing:
		return ctx.Status(202).JSON(response)
	case redis.StatusSuccess:
		return ctx.JSON(response)
	default:
		return ctx.Status(500).JSON(fiber.Map{
			"status": "Error",
			"error":  "Unknown status",
		})
	}
}

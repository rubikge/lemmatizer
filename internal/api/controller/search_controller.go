package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rubikge/lemmatizer/internal/dto"
	"github.com/rubikge/lemmatizer/internal/search"
)

type SearchController struct {
	searchService *search.Service
}

func NewSearchController(s *search.Service) *SearchController {
	return &SearchController{searchService: s}
}

func (c *SearchController) ProcessHandler(ctx fiber.Ctx) error {
	var requestData dto.RequestData
	if err := ctx.Bind().JSON(&requestData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "Error",
			"error":  "Invalid request format",
		})
	}

	result, err := c.searchService.ProcessSearch(&requestData)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Error",
			"error":  "Failed to process search request",
		})
	}

	return ctx.JSON(result)
}

func (c *SearchController) GetResultHandler(ctx fiber.Ctx) error {
	taskID := ctx.Params("taskID")
	result, err := c.searchService.GetResult(taskID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Error",
			"error":  "Failed to get search result",
		})
	}

	switch result.Status {
	case dto.StatusError:
		return ctx.Status(fiber.StatusInternalServerError).JSON(result)
	case dto.StatusProcessing:
		return ctx.Status(102).JSON(result)
	case dto.StatusSuccess:
		return ctx.JSON(result)
	case dto.StatusNotFound:
		return ctx.Status(fiber.StatusNotFound).JSON(result)
	case dto.StatusWrongTaskID:
		return ctx.Status(fiber.StatusNotFound).JSON(result)
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Error",
			"error":  "Unknown status",
		})
	}
}

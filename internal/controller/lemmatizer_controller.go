package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rubikge/lemmatizer/internal/models"
	"github.com/rubikge/lemmatizer/internal/services"
	"github.com/rubikge/lemmatizer/internal/utils"
)

type LemmatizerFiberController struct {
	s *services.LemmatizerService
}

func NewLemmatizerFiberController(s *services.LemmatizerService) *LemmatizerFiberController {
	return &LemmatizerFiberController{s: s}
}

func (c *LemmatizerFiberController) LemmatizeHandler(ctx fiber.Ctx) error {
	var request struct {
		Text string `json:"text"`
	}

	if err := ctx.Bind().JSON(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON",
		})
	}

	lemmas, err := c.s.GetLemmas(request.Text)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"lemmas": lemmas,
	})
}

func (c *LemmatizerFiberController) SearchHandler(ctx fiber.Ctx) error {
	var requestData models.RequestData

	if err := ctx.Bind().JSON(&requestData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON",
		})
	}

	lemmas, err := c.s.GetLemmas(requestData.Message)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	searchProducts, err := utils.GetLemmatizedSearchProduct(&requestData.Product, c.s)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	result := services.GetScore(&lemmas, &searchProducts)
	return ctx.JSON(result)
}

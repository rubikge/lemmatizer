package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rubikge/lemmatizer/internal/services"
)

type LemmatizerFiberController struct {
	s *services.LemmatizerService
}

func NewLemmatizerFiberController(s *services.LemmatizerService) *LemmatizerFiberController {
	return &LemmatizerFiberController{s: s}
}

func (c *LemmatizerFiberController) ProcessText(ctx fiber.Ctx) error {
	var request struct {
		Text string `json:"text"`
	}

	if err := ctx.Bind().JSON(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON",
		})
	}

	lemmas, err := c.s.GetLemmasArray(request.Text)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"lemmas": lemmas,
	})
}

package controller

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rubikge/lemmatizer/internal/lemmatizer"
)

type LemmatizeController struct {
	ls *lemmatizer.LemmatizerService
}

func NewLemmatizeController(ls *lemmatizer.LemmatizerService) *LemmatizeController {
	return &LemmatizeController{ls: ls}
}

func (c *LemmatizeController) LemmatizeHandler(ctx fiber.Ctx) error {
	var request struct {
		Text string `json:"text"`
	}

	if err := ctx.Bind().JSON(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON",
		})
	}

	lemmas, err := c.ls.GetLemmas(request.Text)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(fiber.Map{
		"lemmas": lemmas,
	})
}

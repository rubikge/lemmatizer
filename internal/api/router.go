package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rubikge/lemmatizer/internal/api/controller"
	"github.com/rubikge/lemmatizer/internal/search"
)

func Router(app *fiber.App, s *search.Service) {
	lc := controller.NewLemmatizeController(s.Lemmatizer)
	sc := controller.NewSearchController(s)

	app.Post("/lemmatize", lc.LemmatizeHandler)
	app.Post("/process", sc.ProcessHandler)
	app.Get("/result/:taskID", sc.GetResultHandler)
}

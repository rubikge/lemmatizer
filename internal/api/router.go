package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rubikge/lemmatizer/internal/api/controller"
	"github.com/rubikge/lemmatizer/internal/lemmatizer"
	"github.com/rubikge/lemmatizer/internal/redis"
)

func Router(app *fiber.App, rq *redis.RedisQueue, ls *lemmatizer.LemmatizerService) {
	lc := controller.NewLemmatizeController(ls)
	sc := controller.NewSearchController(rq)

	app.Post("/lemmatize", lc.LemmatizeHandler)
	app.Post("/process", sc.ProcessHandler)
	app.Get("/result/:taskID", sc.GetResultHandler)
}

package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rubikge/lemmatizer/internal/controller"
	"github.com/rubikge/lemmatizer/internal/repository"
	"github.com/rubikge/lemmatizer/internal/services"
)

func main() {
	r := repository.NewMystemRepository()

	s := services.NewLemmatizerService(r)

	c := controller.NewLemmatizerFiberController(s)

	app := fiber.New()
	app.Post("/lemmatize", c.LemmatizeHandler)
	app.Post("/search", c.SearchHandler)
	app.Listen(":3000")
}

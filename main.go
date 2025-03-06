package main

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/rubikge/lemmatizer/internal/api"
	"github.com/rubikge/lemmatizer/internal/search"
)

func main() {
	s, err := search.NewService("mySearch")
	if err != nil {
		fmt.Println(err)
		return
	}

	app := fiber.New()
	api.Router(app, s)

	app.Listen(":3000")
}

package main

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/rubikge/lemmatizer/internal/api"
	"github.com/rubikge/lemmatizer/internal/redis"
	"github.com/rubikge/lemmatizer/internal/search"
)

func main() {
	rq, err := redis.NewRedisQueue()
	if err != nil {
		fmt.Println(err)
		return
	}

	s := search.NewSearchService()

	app := fiber.New()
	api.Router(app, rq, s.Lemmatizer)

	rq.StartWorker("search_worker", s.GetScore)

	app.Listen(":3000")
}

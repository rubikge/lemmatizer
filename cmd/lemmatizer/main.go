package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/rubikge/lemmatizer/internal/controller"
	"github.com/rubikge/lemmatizer/internal/repository"
	"github.com/rubikge/lemmatizer/internal/service"
)

func main() {
	// data, err := os.ReadFile("./testdata/lemmatizer_test_data.txt")
	// if err != nil {
	// 	fmt.Println("Error reading file:", err)
	// 	return
	// }

	r := repository.NewMystemRepository()

	s := service.NewLemmatizerService(r)

	c := controller.NewLemmatizerFiberController(s)

	app := fiber.New()
	app.Post("/lemmatize", c.ProcessText)
	app.Listen(":3000")

}

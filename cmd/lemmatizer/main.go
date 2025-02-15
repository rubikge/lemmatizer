package main

import (
	"fmt"
	"os"

	"github.com/rubikge/lemmatizer/internal/repository"
	"github.com/rubikge/lemmatizer/internal/service"
)

func main() {
	data, err := os.ReadFile("./testdata/lemmatizer_test_data.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	r := repository.NewMystemRepository(data)

	uc := service.NewLemmatizerService(r)

	lemmas, err := uc.ProcessData(string(data))
	if err != nil {
		fmt.Println("Error processing data:", err)
		return
	}

	fmt.Println(lemmas)
}

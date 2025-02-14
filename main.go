package main

import (
	"fmt"
	"lemmatizer-app/mystem"
)

func main() {
	text := "Кошки мурлыкают, когда им хорошо."
	lemmas, err := mystem.Lemmatize(text)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, lemma := range lemmas {
		fmt.Println(lemma)
	}
}

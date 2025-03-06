package positive_check

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type PositiveWords struct {
	Words []string `json:"words"`
}

type Service struct {
	positiveWords map[string]struct{}
}

func NewService() (*Service, error) {
	content, err := os.ReadFile("internal/src/dist/positive/ru.json")
	if err != nil {
		return nil, err
	}

	var words PositiveWords
	if err := json.Unmarshal(content, &words); err != nil {
		return nil, err
	}

	// Create map for O(1) lookup
	positiveWords := make(map[string]struct{}, len(words.Words))
	for _, word := range words.Words {
		positiveWords[word] = struct{}{}
	}

	return &Service{
		positiveWords: positiveWords,
	}, nil
}

func (s *Service) ContainsPositiveWords(message string) bool {
	// Convert to lower case and remove newlines
	message = strings.ToLower(strings.ReplaceAll(message, "\n", " "))

	// Split into words
	words := strings.Fields(message)

	fmt.Printf("Checking message for positive words: %s\n", message)

	// Check each word if it contains any positive word as a prefix
	for _, word := range words {
		for positiveWord := range s.positiveWords {
			if strings.HasPrefix(word, positiveWord) {
				fmt.Printf("Found positive word: %s (matched with: %s)\n", word, positiveWord)
				return true
			}
		}
	}

	fmt.Println("No positive words found in the message")
	return false
}

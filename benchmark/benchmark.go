package benchmark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func generateRandomString(words []string) string {
	rand.Shuffle(len(words), func(i, j int) {
		words[i], words[j] = words[j], words[i]
	})
	return strings.Join(words, " ")
}

func sendRequest(text string) (time.Duration, error) {
	requestBody, err := json.Marshal(map[string]string{
		"text": text,
	})
	if err != nil {
		return 0, err
	}

	buffer := bytes.NewBuffer(requestBody)
	start := time.Now()
	resp, err := http.Post("http://localhost:3000/lemmatize", "application/json", buffer)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	elapsed := time.Since(start)
	return elapsed, nil
}

func recordToFile(f *os.File, str string) {
	if _, err := f.WriteString(str); err != nil {
		log.Fatal(err)
	}
}

func RunTest(words []string, iterations int) {
	times := make([]time.Duration, 0, iterations)

	for i := 0; i < iterations; i++ {
		randomString := generateRandomString(words)
		elapsed, err := sendRequest(randomString)
		if err != nil {
			fmt.Printf("Error sending request: %v\n", err)
		} else {
			fmt.Printf("Test %d: %v\n", i+1, elapsed)
			times = append(times, elapsed)
		}
	}

	var totalTime time.Duration
	minTime := time.Duration(math.MaxInt64)
	maxTime := time.Duration(0)

	for _, t := range times {
		totalTime += t
		if t < minTime {
			minTime = t
		}
		if t > maxTime {
			maxTime = t
		}
	}

	avgTime := totalTime / time.Duration(iterations)

	f, err := os.OpenFile(TestFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	header := fmt.Sprintf("Test passed!\nNumber of words in text: %d\nNumber of tests: %d\n", len(words), iterations)
	recordToFile(f, header)
	stats := fmt.Sprintf("Min time: %v\nMax time: %v\nAvg time: %v\n\n", minTime, maxTime, avgTime)
	recordToFile(f, stats)
}

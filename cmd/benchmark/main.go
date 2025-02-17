package main

import (
	"github.com/rubikge/lemmatizer/benchmark"
)

func main() {
	benchmark.RunTest(benchmark.TestWords, 5)
}

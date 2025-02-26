package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rubikge/lemmatizer/internal/models"
)

func RunSearchTest(jsonData []byte) {
	buffer := bytes.NewBuffer(jsonData)

	resp, err := http.Post("http://127.0.0.1:3000/search", "application/json", buffer)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var responseData models.SearchResult
	if err := json.Unmarshal(body, &responseData); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	fmt.Printf("Response Struct: %+v\n", responseData)
}

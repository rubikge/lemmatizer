package tests

import (
	"bytes"
	"fmt"
	"net/http"
)

func RunSearchTest(jsonData []byte) {
	buffer := bytes.NewBuffer(jsonData)

	resp, err := http.Post("http://localhost:3000/search", "application/json", buffer)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp)

}

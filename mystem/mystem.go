package mystem

import (
	"encoding/json"
	"os/exec"
	"strings"
)

type LemmaResponse struct {
	Lemma string `json:"lemma"`
}

func Lemmatize(text string) ([]LemmaResponse, error) {
	cmd := exec.Command("./mystem/mystem", "-l", "-n", "--json")
	cmd.Stdin = strings.NewReader(text)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var lemmas []LemmaResponse
	err = json.Unmarshal(out, &lemmas)
	if err != nil {
		return nil, err
	}
	return lemmas, nil
}

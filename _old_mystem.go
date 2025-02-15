package mystem

import (
	"bufio"
	"encoding/json"
	"os/exec"
	"strings"
)

type Analysis struct {
	Lex string `json:"lex"`
	Gr  string `json:"gr"`
}

type AnalizedWord struct {
	Analysis []Analysis `json:"analysis"`
	Text     string     `json:"text"`
}

func Lemmatize(text string) ([]string, error) {
	cmd := exec.Command("./mystem/mystem", "-dnig", "--format", "json")
	cmd.Stdin = strings.NewReader(text)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	needePrefixes := []string{"A", "ADV", "COM", "S", "V"}
	var lemmas []string

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		var word AnalizedWord
		if err := json.Unmarshal([]byte(line), &word); err != nil {
			continue
		}

		if len(word.Analysis) == 0 {
			continue
		}

		analysis := word.Analysis[0]

		for _, prefix := range needePrefixes {
			if strings.HasPrefix(analysis.Gr, prefix+"=") ||
				strings.HasPrefix(analysis.Gr, prefix+",") {
				lemmas = append(lemmas, analysis.Lex)
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return lemmas, nil
}

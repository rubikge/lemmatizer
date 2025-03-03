package mystem

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type MystemRepository struct{}

func NewMystemRepository() *MystemRepository {
	return &MystemRepository{}
}

func (r *MystemRepository) GetAnalysis(text string) ([]AnalizedWord, error) {
	myStem := exec.Command(MystemExecPath, MystemFlags, "--format", MystemFormat)
	myStem.Stdin = strings.NewReader(text)

	stdout, err := myStem.StdoutPipe()
	if err != nil {
		fmt.Printf("error getting stdout pipe: %v", err)
		return nil, err
	}

	if err := myStem.Start(); err != nil {
		fmt.Printf("error starting mystem: %v", err)
		return nil, err
	}

	scanner := bufio.NewScanner(stdout)
	var analysis []AnalizedWord

	for scanner.Scan() {
		line := scanner.Text()
		word, err := parseLine(line)

		if err != nil {
			fmt.Printf("error parsing line: %v", err)
			continue
		}

		analysis = append(analysis, word)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error scanning output: %v", err)
	}

	if err := myStem.Wait(); err != nil {
		fmt.Printf("error waiting for mystem: %v", err)
	}

	return analysis, nil
}

func parseLine(line string) (AnalizedWord, error) {
	var word AnalizedWord
	if err := json.Unmarshal([]byte(line), &word); err != nil {
		return word, err
	}

	return word, nil
}

package repository

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/rubikge/lemmatizer/internal/models"
	"github.com/rubikge/lemmatizer/internal/mystem"
)

type MystemRepository struct {
	myStem *exec.Cmd
}

func NewMystemRepository() *MystemRepository {
	myStem := exec.Command(mystem.MystemExecPath, mystem.MystemFlags, "--format", mystem.MystemFormat)
	return &MystemRepository{
		myStem: myStem,
	}
}

func (r *MystemRepository) GetAnalysis(text string) ([]models.AnalizedWord, error) {
	r.myStem.Stdin = strings.NewReader(text)

	stdout, err := r.myStem.StdoutPipe()
	if err != nil {
		fmt.Printf("error getting stdout pipe: %v", err)
		return nil, err
	}

	if err := r.myStem.Start(); err != nil {
		fmt.Printf("error starting mystem: %v", err)
		return nil, err
	}

	scanner := bufio.NewScanner(stdout)
	var analysis []models.AnalizedWord

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

	if err := r.myStem.Wait(); err != nil {
		fmt.Printf("error waiting for mystem: %v", err)
	}

	return analysis, nil
}

func parseLine(line string) (models.AnalizedWord, error) {
	var word models.AnalizedWord
	if err := json.Unmarshal([]byte(line), &word); err != nil {
		return word, err
	}

	return word, nil
}

package repository

import (
	"bufio"
	"encoding/json"
	"log"
	"os/exec"

	"github.com/rubikge/lemmatizer/internal/model"
	"github.com/rubikge/lemmatizer/internal/mystem"
)

type MystemRepository struct {
}

func NewMystemRepository() *MystemRepository {
	return &MystemRepository{}
}

func (r *MystemRepository) GetDataStream(text string) (<-chan model.AnalizedWord, error) {
	wordChan := make(chan model.AnalizedWord)

	cmd := exec.Command(mystem.MystemExecPath, mystem.MystemFlags, "--format", mystem.MystemFormat)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Printf("error getting stdin pipe: %v", err)
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("error getting stdout pipe: %v", err)
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		log.Printf("error starting mystem: %v", err)
		return nil, err
	}

	go func() {
		defer close(wordChan)

		go func() {
			defer stdin.Close()
			stdin.Write([]byte(text))
		}()

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			word, err := parseLine(line)
			if err != nil {
				log.Printf("error parsing line: %v", err)
				continue
			}
			wordChan <- word
		}

		if err := scanner.Err(); err != nil {
			log.Printf("error scanning output: %v", err)
		}

		if err := cmd.Wait(); err != nil {
			log.Printf("error waiting for mystem: %v", err)
		}
	}()

	return wordChan, nil
}

func parseLine(line string) (model.AnalizedWord, error) {
	var word model.AnalizedWord
	if err := json.Unmarshal([]byte(line), &word); err != nil {
		return word, err
	}

	return word, nil
}

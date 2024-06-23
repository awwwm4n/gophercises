package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type Quiz struct {
	question, answer string
}

func main() {

	fileName := flag.String("csv", "problems.csv", "a csv file in the format of 'question, answer'")
	timeLimit := flag.Int("time_limit", 30, "time limit for all questions in seconds")

	flag.Parse()

	quizzes, err := readCsv(*fileName)

	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	timer := time.NewTimer(time.Second * time.Duration(*timeLimit))

	correctAnswers := 0

	go func() {
		for _, quiz := range quizzes {
			fmt.Printf("%s: ", quiz.question)
			var userAnswer string
			fmt.Scan(&userAnswer)
			if strings.TrimSpace(userAnswer) == quiz.answer {
				correctAnswers++
			}
		}
		done <- true
	}()
	select {
	case <-done:
		timer.Stop()
	case <-timer.C:
		fmt.Println("time's up")
	}

	fmt.Printf("You scored: %v/%v", correctAnswers, len(quizzes))

}

func readCsv(fileName string) ([]Quiz, error) {

	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("error while opening file: %v", err)
	}

	csvReader := csv.NewReader(file)

	quizzes := []Quiz{}

	for {

		record, err := csvReader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("error while reading csv record: %v", err)
		}

		if len(record) != 2 {
			return nil, fmt.Errorf("invalid no of fiels in csv for record: %v", record)
		}

		quizzes = append(quizzes, Quiz{record[0], record[1]})

	}

	return quizzes, nil
}

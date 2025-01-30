package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

// Quiz represents the quiz with questions and answers.
type Quiz struct {
	Questions [][]string
	Correct   int
}

// NewQuiz initializes a new Quiz from a CSV file.
func NewQuiz(fileName string) (*Quiz, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %s", fileName)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	questions, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading questions: %v", err)
	}

	return &Quiz{Questions: questions}, nil
}

// Run starts the quiz and handles timing and user input.
func (q *Quiz) Run(timeLimit int) {
	fmt.Printf("There are %d questions\n", len(q.Questions))
	fmt.Print("Press return key when ready\n")
	fmt.Scanln()

	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	defer timer.Stop()

	startTime := time.Now()

	for i, question := range q.Questions {
		fmt.Printf("Question %d - What is: %v = ", i+1, strings.TrimSpace(question[0]))

		answerCh := make(chan string)
		go func() {
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			answerCh <- strings.TrimSpace(scanner.Text())
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nTime up! You scored %d out of %d\n", q.Correct, len(q.Questions))
			return
		case ans := <-answerCh:
			if ans == question[1] {
				q.Correct++
			}
		}
	}

	timer.Stop()
	fmt.Printf("You scored %d out of %d in %.2f seconds\n", q.Correct, len(q.Questions), time.Since(startTime).Seconds())
}

// exit prints a message and exits the program with the given code.
func exit(msg string, code int) {
	fmt.Println(msg)
	os.Exit(code)
}

func main() {
	timeLimit := flag.Int("timeLimit", 20, "quiz timer")
	fileName := flag.String("fileName", "questions.csv", "a quiz file in csv format")
	flag.Parse()

	quiz, err := NewQuiz(*fileName)
	if err != nil {
		exit(err.Error(), 1)
	}

	quiz.Run(*timeLimit)
}

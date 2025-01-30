package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
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

// Config represents quiz configuration from JSON.
type Config struct {
	CSVFile   string `json:"csv_file"`
	Delimiter string `json:"delimiter"`
	TimeLimit int    `json:"time_limit"`
}

// LoadConfig reads the quiz settings from a JSON file.
func LoadConfig(configFile string) (*Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file '%s': %v", configFile, err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config file '%s': %v", configFile, err)
	}

	// Validate delimiter
	if len(config.Delimiter) != 1 {
		return nil, fmt.Errorf("invalid delimiter '%s' in config file", config.Delimiter)
	}

	return &config, nil
}

// NewQuiz initializes a new Quiz from a CSV file.
func NewQuiz(fileName string, delimiter rune) (*Quiz, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file '%s': %v", fileName, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = delimiter // Set custom delimiter

	questions, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading questions from '%s': %v", fileName, err)
	}

	if len(questions) == 0 {
		return nil, fmt.Errorf("CSV file '%s' is empty", fileName)
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

	fmt.Printf("You scored %d out of %d in %.2f seconds\n", q.Correct, len(q.Questions), time.Since(startTime).Seconds())
}

// exit prints a message and exits the program with the given code.
func exit(msg string, code int) {
	fmt.Println(msg)
	os.Exit(code)
}

func main() {
	// Command-line flags
	timeLimit := flag.Int("timeLimit", 20, "quiz timer in seconds")
	fileName := flag.String("fileName", "questions.csv", "CSV file with quiz questions")
	delimiter := flag.String("delimiter", ",", "CSV delimiter character")
	configFile := flag.String("config", "", "optional config JSON file")
	flag.Parse()

	// Load configuration from JSON if provided
	if *configFile != "" {
		config, err := LoadConfig(*configFile)
		if err != nil {
			exit(err.Error(), 1)
		}
		*fileName = config.CSVFile
		*delimiter = config.Delimiter
		*timeLimit = config.TimeLimit
	}

	// Validate delimiter
	if len(*delimiter) != 1 {
		exit("Invalid delimiter: must be a single character", 1)
	}

	// Load the quiz
	quiz, err := NewQuiz(*fileName, rune((*delimiter)[0]))
	if err != nil {
		exit(err.Error(), 1)
	}

	// Start the quiz
	quiz.Run(*timeLimit)
}

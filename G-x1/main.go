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

func main() {

	timeLimit := flag.Int("timeLimit", 20, "quiz timer")
	fileName := flag.String("fileName", "questions.csv", "a quiz file in csv format")
	flag.Parse()

	file, err := os.Open(*fileName)

	if err != nil {
		exit(fmt.Sprintf("Failed to open CSV file: %s\n", *fileName), 1)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	questions, err := reader.ReadAll()

	if err != nil {
		exit("Error reading questions", 1)
	}

	//keep track of number of correct answers
	correct := 0
	//display total number of questions
	fmt.Printf("There are %d questions\n", len(questions))
	fmt.Print("Press return key when ready\n")
	fmt.Scanln()

	// initiate timer
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	defer timer.Stop()

	//run timer
	go func() {
		<-timer.C
		exit(fmt.Sprintf("\nTime up! You scored %d out of %d", correct, len(questions)), 0)

	}()
	startTime := time.Now()
	for i, question := range questions {
		// ask the next question
		fmt.Printf("Question %d - What is:  %v = ", i+1, strings.TrimSpace(question[0]))

		//wait for answer
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		ans := strings.TrimSpace(scanner.Text())

		//compare input answer to actual answer
		if ans == question[1] {
			//increment correct if correct answer given
			correct += 1
		}
	}
	//stop timer if all questions have been answered in the time
	timer.Stop()

	//return correct answers and elapsed time
	fmt.Printf("You scored %d out of %d in %.2f seconds\n", correct, len(questions), (time.Since(startTime)).Seconds())
}

func exit(msg string, code int) {
	fmt.Println(msg)
	os.Exit(code)
}

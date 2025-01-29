package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {

	timeLimit := flag.Int("timeLimit", 20, "quiz timer")
	fileName := flag.String("fileName", "questions.csv", "a quiz file in csv format")
	flag.Parse()

	file, err := os.Open(*fileName)

	if err != nil {
		log.Fatal("Error while reading file", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	questions, err := reader.ReadAll()

	if err != nil {
		fmt.Println("Error reading questions")
	}

	//keep track of number of correct answers
	correct := 0
	//display total number of questions
	fmt.Printf("There are %v questions\n", len(questions))
	fmt.Print("Press return key when ready\n")
	fmt.Scanln()

	// initiate timer
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	defer timer.Stop()

	//run timer
	go func() {
		<-timer.C
		fmt.Printf("\nTime up! You scored %v out of %v", correct, len(questions))
		//exit program
		os.Exit(0)
	}()
	startTime := time.Now()
	for i, question := range questions {
		// ask the next question
		fmt.Printf("Question %v - What is:  %v\n", i+1, question[0])

		//wait for answer
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		ans := scanner.Text()

		//compare input answer to actual answer
		if ans == question[1] {
			//increment correct if correct answer given
			correct += 1
		}
	}
	//stop timer if all questions have been answered in the time
	timer.Stop()

	//return correct answers and elapsed time
	fmt.Printf("You scored %v out of %v in %.2f seconds\n", correct, len(questions), (time.Since(startTime)).Seconds())
}

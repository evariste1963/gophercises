package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func main() {

	file, err := os.Open("questions.csv")

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

	for i, question := range questions {
		ans := ""
		// ask the next question
		fmt.Printf("Question %v - Calculate %v\n", i+1, question[0])
		//wait for answer
		fmt.Scanln(&ans)
		//compare input answer to actual answer
		if ans == question[1] {
			//increment correct if correct answer given
			correct += 1
		}
	}
	//return correct answers no etc
	fmt.Printf("You scored %v out of %v\n", correct, len(questions))
}

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

	//total number of questions
	fmt.Printf("There are %v questions\n", len(questions))
	// END of TEST
	correct := 0
	for i, eachquestion := range questions {
		ans := ""
		// ask the next question, inc question number
		fmt.Printf("Question %v - Calculate %v\n", i+1, eachquestion[0])
		//wait for answer
		fmt.Scanln(&ans)
		//compare input answer to actual answer
		if ans == eachquestion[1] {
			//increment correct if correct answer given
			correct += 1
		}
	}
	//return correct answers no etc
	fmt.Printf("You scored %v out of %v\n", correct, len(questions))

}

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

	// reading 1st row, 1st column
	question_1 := questions[0][0]

	fmt.Println(question_1)
	// END of TEST

	for _, eachquestion := range questions {
		fmt.Println(eachquestion)
	}

}

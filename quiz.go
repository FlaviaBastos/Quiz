package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
)

// Line holds each line from file
type Line struct {
	Question string
	Answer   string
}

func main() {
	var filename string
	flag.StringVar(&filename, "csv", "questions.csv", "A csv file with one question and answer per line")
	flag.Parse()

	// Open file
	csvfile, err := os.Open(filename)
	if err != nil {
		fmt.Printf("cannot open file %s. Error: %v", filename, err)
		os.Exit(1)
	}

	// Reads file
	r := csv.NewReader(csvfile)

	var quiz []Line
	// Saves file content into struct
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Cannot read file", err)
		}

		quiz = append(quiz, Line{
			Question: line[0],
			Answer:   line[1],
		})
	}

	var score int
	fmt.Println("Ready? Go!")
	for _, question := range quiz {
		var input string
		fmt.Printf("Question: %v\n", question.Question)
		fmt.Print("Type your answer: ")
		_, err := fmt.Scanln(&input)

		if err != nil {
			fmt.Println("Cannot read your answer")
		}

		if input == question.Answer {
			fmt.Println("Correct")
			score++
		} else {
			fmt.Println("Wrong")
		}

	}

	fmt.Printf("Your score is: %v out of %v\n", score, len(quiz))
}

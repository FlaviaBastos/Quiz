package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

// Line holds each line from file
type Line struct {
	Question string
	Answer   string
}

func main() {
	var filename string
	var timer string
	flag.StringVar(&filename, "csv", "questions.csv", "A csv file with one question and answer per line")
	flag.StringVar(&timer, "time", "30", "Quiz time limit in seconds")
	flag.Parse()

	// Open file
	csvfile, err := os.Open(filename)
	if err != nil {
		fmt.Printf("cannot open file %s. Error: %v", filename, err)
		os.Exit(1)
	}

	// Reads file
	r := csv.NewReader(csvfile)

	lines, err := r.ReadAll()
	if err == io.EOF {
		os.Exit(1)
	}
	if err != nil {
		fmt.Println("cannot read file", err)
		os.Exit(1)
	}

	// creates quiz from file
	quiz := parse(lines)

	// determine time cap
	t, err := strconv.Atoi(timer)
	if err != nil {
		fmt.Println("this is not a valid timeout value", err)
	}
	timelimit := time.Duration(t) * time.Second

	// get ready to start game
	fmt.Println("Ready? Hit <ENTER> to start!")
	fmt.Scanln()

	// start game
	c := make(chan int)
	go runQuiz(quiz, c)
	timeout := time.After(timelimit)

	// wait for either quiz to finish or time runs out
	for {
		select {
		case total := <-c:
			fmt.Printf("Your score is: %v out of %v\n", total, len(quiz))
			return
		case <-timeout:
			fmt.Println("Time is up!")
			return
		}
	}
}

func parse(lines [][]string) []Line {
	// Saves file content into struct
	ques := make([]Line, len(lines)) // using make instead of append because the length is known, so no need to adjust size everytime
	for i, line := range lines {
		ques[i] = Line{
			Question: line[0],
			Answer:   line[1],
		}
	}
	return ques
}

func runQuiz(quiz []Line, c chan int) {
	var score int
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
	c <- score
}

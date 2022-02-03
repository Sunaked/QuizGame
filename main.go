package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {

	csvFilename := flag.String("csv", "../input.csv", "a csv file in the format of 'question, answer")
	amountOfQuestions := flag.Int("amountOfQuestions", 0, "number of questions to display")
	needToShuffle := flag.Bool("shuffle", false, "needToShuffles the question order")
	timeForQuestions := flag.Int("time", 30, "time to answer questions, after exceeding time test will be over")

	flag.Parse()

	_ = needToShuffle
	_ = amountOfQuestions
	_ = timeForQuestions

	//Geting problems
	/////////////////////////////
	file, err := os.Open(*csvFilename)
	if err != nil {
		fmt.Printf("Error opening csv file: %s\n", *csvFilename)
		os.Exit(1)
	}
	f := csv.NewReader(file)
	lines, err := f.ReadAll()
	if err != nil {
		fmt.Printf("Error reading csv: %s\n", err)
		os.Exit(1)
	}
	problems, err := getProblems(lines)
	if err != nil {
		fmt.Printf("Error getting problems: %s\n", err)
		os.Exit(1)
	}
	/////////////////////////////
	if *amountOfQuestions == 0 {
		*amountOfQuestions = len(problems) - 1
	}
	var correctAnswers int = 0
	var amountOfAnswers int = 0
	startGame()

	timer := time.NewTimer(time.Second * time.Duration(*timeForQuestions))
	ch1 := make(chan string)
GameLoop:
	for ind := 0; ind < *amountOfQuestions; ind++ {
		rand.Seed(int64(time.Now().UnixNano()))
		var problem problem
		if *needToShuffle {
			problemID := rand.Intn(len(problems) - 1)
			problems = append(problems[:problemID], problems[problemID+1:]...)
			problem = problems[problemID]
		} else {
			problems = append(problems[:0], problems[1:]...)
			fmt.Println(problems)
			problem = problems[0]
		}
		fmt.Printf("Question is: %s, your answer = ", problem.question)
		amountOfAnswers++
		var answer string
		go func() {
			var tmp string
			fmt.Scanf("%s\n", &tmp)
			ch1 <- tmp
		}()
		select {
		case <-timer.C:
			break GameLoop
		case answer = <-ch1:
			if answer == problem.answer {
				fmt.Print("Rigth answer! Good\n")
				correctAnswers++
			} else {
				fmt.Printf("Wrong! Bad: Right answer is %s, whereas your answer is %s\n", problem.answer, answer)
			}
		}
	}
	gameOver(correctAnswers, amountOfAnswers)
}

func startGame() {
	fmt.Printf("\nPress any key to begin")
	var tmp string
	fmt.Scanf("%s\n", tmp)
}

func gameOver(correctAnswers int, totalAnswers int) {
	fmt.Printf("\n\nGame is over, Your score is %d out of %d\n", correctAnswers, totalAnswers)
}

func getProblems(lines [][]string) ([]problem, error) {
	problems := make([]problem, len(lines))
	for ind, line := range lines {
		problems[ind] = problem{
			question: strings.TrimSpace(line[0]),
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return problems, nil
}

type problem struct {
	question string
	answer   string
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	Win  = 6
	Lose = 0
	Draw = 3
)

const (
	Rock    = 1
	Paper   = 2
	Scissor = 3
)

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	total := 0
	for scanner.Scan() {
		line := scanner.Text()
		round := strings.Split(line, " ")
		fmt.Println(round, score(round[0], round[1]))
		total += score(round[0], round[1])
	}

	fmt.Println(total)
}

// part 2
func lose(s string) int {
	switch s {
	case "A":
		return Scissor
	case "B":
		return Rock
	case "C":
		return Paper
	}
	panic("lose")
}

func draw(s string) int {
	switch s {
	case "A":
		return Rock
	case "B":
		return Paper
	case "C":
		return Scissor
	}
	panic("draw")
}

func win(s string) int {
	switch s {
	case "A":
		return Paper
	case "B":
		return Scissor
	case "C":
		return Rock
	}
	panic("win")
}

func score(a, b string) int {
	switch b {
	case "X":
		return Lose + lose(a)
	case "Y":
		return Draw + draw(a)
	case "Z":
		return Win + win(a)
	}

	panic("score")
}

/*
part 1
func outcome(a, b string) int {
	switch a + b {
	case "AX":
		return Draw
	case "AY":
		return Win
	case "AZ":
		return Lose
	case "BX":
		return Lose
	case "BY":
		return Draw
	case "BZ":
		return Win
	case "CX":
		return Win
	case "CY":
		return Lose
	case "CZ":
		return Draw
	}

	panic("not a valid case: " + a + " - " + b)
}

func shape(s string) int {
	switch s {
	case "X":
		return 1
	case "Y":
		return 2
	case "Z":
		return 3
	}

	panic("not a valid shape: " + s)
}
*/

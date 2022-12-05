package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Stack struct {
	items []string
}

func (s *Stack) Pop() string {
	item := s.Peek()
	s.items = s.items[:len(s.items)-1]
	return item
}

func (s *Stack) Push(item string) {
	s.items = append(s.items, item)
}

func (s *Stack) PopN(n int) []string {
	items := s.items[len(s.items)-n:]
	s.items = s.items[:len(s.items)-n]
	return items
}

func (s *Stack) PushN(items []string) {
	s.items = append(s.items, items...)
}

func (s *Stack) Peek() string {
	return s.items[len(s.items)-1]
}

type Crates struct {
	stacks []Stack
}

func (c *Crates) MovePart1(mv Move) {
	for i := 0; i < mv.count; i++ {
		crate := c.stacks[mv.from].Pop()
		c.stacks[mv.to].Push(crate)
	}
}

func (c *Crates) MovePart2(mv Move) {
	crates := c.stacks[mv.from].PopN(mv.count)
	c.stacks[mv.to].PushN(crates)
}

func ParseCrates(scanner *bufio.Scanner) Crates {
	var lines []string

	var stacks []Stack
	for scanner.Scan() {
		line := scanner.Text()

		// find the line with the number of stacks
		if strings.HasPrefix(line, " 1 ") {
			lastCol, err := strconv.Atoi(string(line[len(line)-2]))
			if err != nil {
				panic(err)
			}
			stacks = make([]Stack, lastCol)
			break
		}

		lines = append(lines, line)
	}

	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		stack := 0
		for j := 0; j < len(line); {
			if line[j] == '[' {
				j += 1
				stacks[stack].Push(string(line[j]))
				j += 1
				if line[j] != ']' {
					panic(line)
				}
				j += 2
			} else {
				j += 4
			}
			stack += 1
		}
	}

	return Crates{
		stacks: stacks,
	}
}

type Move struct {
	count int
	from  int
	to    int
}

func ParseMoves(scanner *bufio.Scanner) []Move {
	var moves []Move
	var count, from, to int
	for scanner.Scan() {
		n, err := fmt.Sscanf(scanner.Text(), "move %d from %d to %d", &count, &from, &to)
		if err != nil {
			panic(err)
		}

		if n != 3 {
			panic(err)
		}

		moves = append(moves, Move{count, from - 1, to - 1})
	}
	return moves
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	crates := ParseCrates(scanner)
	fmt.Println(crates)

	scanner.Scan()
	if scanner.Text() != "" {
		panic(scanner.Text())
	}
	moves := ParseMoves(scanner)
	fmt.Println(moves)

	// part1(crates, moves)
	part2(crates, moves)
}

func part1(crates Crates, moves []Move) {
	for _, mv := range moves {
		crates.MovePart1(mv)
		fmt.Println(mv, crates)
	}

	var out string
	for _, stack := range crates.stacks {
		out += stack.Peek()
	}

	fmt.Println(out)
}

func part2(crates Crates, moves []Move) {
	for _, mv := range moves {
		crates.MovePart2(mv)
		fmt.Println(mv, crates)
	}

	var out string
	for _, stack := range crates.stacks {
		out += stack.Peek()
	}

	fmt.Println(out)
}

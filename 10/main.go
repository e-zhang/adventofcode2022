package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var registerX int

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	registerX = 1

	crt := &CRT{}
	crt.Init()

	var cmd *Command
	ctr := 0
	signal := 0
	for cycles := 1; ; cycles += 1 {
		if cmd == nil {
			if !scanner.Scan() {
				break
			}
			line := scanner.Text()
			cmd = Parse(line)
		}

		if cycles == 20 || ctr == 40 {
			fmt.Println(cycles, registerX, cycles*registerX)
			signal += cycles * registerX
			ctr = 0
		}

		crt.Draw(cycles - 1)

		cmd.Tick()
		if cmd.IsFinished() {
			cmd = nil
		}
		ctr += 1
	}

	fmt.Println(signal)
	crt.Print()
}

type Command struct {
	cycles       int
	onCompletion func()
}

func Parse(line string) *Command {
	tokens := strings.Split(line, " ")
	switch tokens[0] {
	case "noop":
		return &Command{1, nil}
	case "addx":
		val, err := strconv.Atoi(tokens[1])
		if err != nil {
			panic(err)
		}
		return &Command{2, func() { registerX += val }}
	}
	panic(line)
}

func (c *Command) Tick() {
	c.cycles -= 1
	if c.cycles == 0 && c.onCompletion != nil {
		c.onCompletion()
	}
}

func (c *Command) IsFinished() bool {
	return c.cycles == 0
}

type CRT struct {
	pixels [][]string
}

func (crt *CRT) Init() {
	crt.pixels = make([][]string, 6)
	for i := range crt.pixels {
		crt.pixels[i] = make([]string, 40)
	}
}

func (crt *CRT) Draw(pos int) {
	row := pos / 40
	col := pos % 40

	if col >= registerX-1 && col <= registerX+1 {
		crt.pixels[row][col] = "#"
	} else {
		crt.pixels[row][col] = "."
	}
}

func (crt *CRT) Print() {
	fmt.Println()
	for _, row := range crt.pixels {
		for _, pixel := range row {
			fmt.Printf(pixel)
		}
		fmt.Println()
	}
	fmt.Println()
}

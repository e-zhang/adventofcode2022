package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	RIGHT = "R"
	LEFT  = "L"
	UP    = "U"
	DOWN  = "D"
)

type Motion struct {
	dir   string
	steps int
}

type Point struct {
	row int
	col int
}

func (p *Point) Move(dir string) {
	switch dir {
	case RIGHT:
		p.row += 1
	case LEFT:
		p.row -= 1
	case UP:
		p.col += 1
	case DOWN:
		p.col -= 1
	}
}

func (p *Point) Follow(other *Point) {
	rowDiff := other.row - p.row
	colDiff := other.col - p.col

	sign := func(num int) int {
		switch {
		case num == 0:
			return 0
		case num < 0:
			return -1
		case num > 0:
			return 1
		}

		panic(num)
	}

	switch {
	case -1 <= rowDiff && rowDiff <= 1 && -1 <= colDiff && colDiff <= 1:
		// no-op
	case rowDiff == 0:
		p.col += sign(colDiff)
	case colDiff == 0:
		p.row += sign(rowDiff)
	default:
		p.row += sign(rowDiff)
		p.col += sign(colDiff)
	}
}

func main() {
	f, err := os.Open("test")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	var motions []Motion
	for scanner.Scan() {
		line := scanner.Text()
		var dir string
		var steps int
		fmt.Sscanf(line, "%s %d", &dir, &steps)
		motions = append(motions, Motion{dir, steps})
	}
	part1(motions)
	part2(motions)
}

func part1(motions []Motion) {
	var head Point
	var tail Point

	visit := make(map[Point]struct{})
	visit[tail] = struct{}{}

	for _, m := range motions {
		for i := 0; i < m.steps; i++ {
			head.Move(m.dir)
			tail.Follow(&head)
			visit[tail] = struct{}{}
		}
	}

	fmt.Println(len(visit))
}

func part2(motions []Motion) {
	rope := make([]Point, 10)

	visit := make(map[Point]struct{})
	visit[rope[len(rope)-1]] = struct{}{}

	for _, m := range motions {
		for i := 0; i < m.steps; i++ {
			rope[0].Move(m.dir)
			for k := 1; k < len(rope); k++ {
				rope[k].Follow(&rope[k-1])
			}
			visit[rope[len(rope)-1]] = struct{}{}
		}
	}

	fmt.Println(len(visit))
}

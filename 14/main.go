package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func sign(x int) int {
	if x < 0 {
		return -1
	}
	if x > 0 {
		return 1
	}
	return 0
}

type Coord struct {
	x int
	y int
}

func (c *Coord) PointsTo(other Coord) []Coord {
	dx := other.x - c.x
	dy := other.y - c.y

	var path []Coord

	if dx == 0 {
		for i := 0; i <= sign(dy)*dy; i++ {
			path = append(path, Coord{c.x, c.y + sign(dy)*i})
		}
	}

	if dy == 0 {
		for i := 0; i <= sign(dx)*dx; i++ {
			path = append(path, Coord{c.x + sign(dx)*i, c.y})
		}
	}

	return path
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	rocks := make(map[Coord]struct{})
	var max int

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		points := strings.Split(line, "->")

		var prev *Coord
		for _, p := range points {
			tokens := strings.Split(strings.TrimSpace(p), ",")
			x, err := strconv.Atoi(tokens[0])
			if err != nil {
				panic(err)
			}
			y, err := strconv.Atoi(tokens[1])
			if err != nil {
				panic(err)
			}

			curr := Coord{x, y}
			if prev != nil {
				for _, c := range prev.PointsTo(curr) {
					rocks[c] = struct{}{}
				}
			}

			if curr.y > max {
				max = curr.y
			}
			prev = &curr
		}
	}

	source := Coord{500, 0}
	sands := make(map[Coord]struct{})
	floor := max + 2

	for {
		sand := Fall(source, floor, rocks, sands)
		sands[sand] = struct{}{}
		if sand == source {
			break
		}

	}

	fmt.Println(len(sands))
}

func Fall(source Coord, floor int, rocks, sands map[Coord]struct{}) Coord {
	sand := source
	for {
		var next *Coord
		for _, delta := range []Coord{
			{0, 1},
			{-1, 1},
			{1, 1},
		} {
			check := Coord{sand.x + delta.x, sand.y + delta.y}
			_, isRock := rocks[check]
			_, isSand := sands[check]
			isFloor := check.y == floor

			if !isRock && !isSand && !isFloor {
				next = &check
				break
			}
		}

		if next == nil {
			break
		}
		sand = *next
	}

	return sand

	return sand
}

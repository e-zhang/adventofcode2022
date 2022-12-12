package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Coord struct {
	row int
	col int
}

type Square struct {
	elevation rune
}

func (s *Square) IsReachableFrom(from *Square) bool {
	return int(s.elevation-from.elevation) <= 1
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	var heightmap [][]*Square
	var start, end Coord

	var starts []Coord

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		row := make([]*Square, len(line))
		for i, h := range line {
			coord := Coord{len(heightmap), i}
			switch h {
			case 'S':
				start = coord
				row[i] = &Square{'a'}
				starts = append(starts, start)
			case 'E':
				end = coord
				row[i] = &Square{'z'}
			default:
				row[i] = &Square{h}
				if h == 'a' {
					starts = append(starts, coord)
				}
			}
		}

		heightmap = append(heightmap, row)
	}

	fmt.Println(shortestPath(heightmap, start, end))

	var paths []int
	for _, s := range starts {
		steps := shortestPath(heightmap, s, end)
		if steps != 0 {
			paths = append(paths, steps)
		}
	}
	sort.Ints(paths)
	fmt.Println(paths[0])
}

func shortestPath(heightmap [][]*Square, start, end Coord) int {
	distances := make([][]int, len(heightmap))
	visited := make([][]bool, len(heightmap))
	for r := range heightmap {
		distances[r] = make([]int, len(heightmap[r]))
		visited[r] = make([]bool, len(heightmap[r]))
	}

	queue := []Coord{start}
	for len(queue) > 0 {
		curr := queue[0]
		if curr == end {
			break
		}
		queue = queue[1:]

		for _, point := range neighbors(heightmap, curr) {
			if !visited[point.row][point.col] {
				visited[point.row][point.col] = true
				distances[point.row][point.col] = distances[curr.row][curr.col] + 1
				queue = append(queue, point)
			}
		}
	}

	return distances[end.row][end.col]
}

func neighbors(heightmap [][]*Square, point Coord) []Coord {
	points := []Coord{}

	for _, delta := range [][]int{
		{0, 1},
		{1, 0},
		{0, -1},
		{-1, 0},
	} {
		r, c := point.row+delta[0], point.col+delta[1]
		if r < 0 || c < 0 || r >= len(heightmap) || c >= len(heightmap[0]) {
			continue
		}

		if heightmap[r][c].IsReachableFrom(heightmap[point.row][point.col]) {
			points = append(points, Coord{r, c})
		}
	}
	return points
}

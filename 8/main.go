package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const (
	LEFT = iota
	RIGHT
	UP
	DOWN
)

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	var grid [][]int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		row := make([]int, len(line))
		for i, c := range line {
			t, err := strconv.Atoi(string(c))
			if err != nil {
				panic(err)
			}

			row[i] = t
		}

		grid = append(grid, row)
	}

	visible := 0
	maxScenic := 0
	scenicLoc := []int{-1, -1}
	for r := range grid {
		for c := range grid[r] {
			if r == 0 || c == 0 || r == len(grid)-1 || c == len(grid[r])-1 {
				visible += 1
				continue
			}

			if isVisible(grid, r, c) {
				visible += 1
			}

			if scenicScore := scenic(grid, r, c); scenicScore > maxScenic {
				maxScenic = scenicScore
				scenicLoc = []int{r, c}
			}
		}
	}
	fmt.Println(visible)
	fmt.Println(scenicLoc, maxScenic)
}

func check(grid [][]int, height, row, col, dir int) bool {
	if row < 0 || col < 0 || row >= len(grid) || col >= len(grid[row]) {
		return true
	}

	switch dir {
	case LEFT:
		return grid[row][col] < height && check(grid, height, row-1, col, dir)
	case RIGHT:
		return grid[row][col] < height && check(grid, height, row+1, col, dir)
	case UP:
		return grid[row][col] < height && check(grid, height, row, col-1, dir)
	case DOWN:
		return grid[row][col] < height && check(grid, height, row, col+1, dir)
	}

	panic(dir)
}

func isVisible(grid [][]int, row, col int) bool {
	height := grid[row][col]

	return check(grid, height, row-1, col, LEFT) ||
		check(grid, height, row+1, col, RIGHT) ||
		check(grid, height, row, col-1, UP) ||
		check(grid, height, row, col+1, DOWN)
}

func scenic(grid [][]int, row, col int) int {
	height := grid[row][col]

	return score(grid, height, row-1, col, LEFT) *
		score(grid, height, row+1, col, RIGHT) *
		score(grid, height, row, col-1, UP) *
		score(grid, height, row, col+1, DOWN)
}

func score(grid [][]int, height, row, col, dir int) int {
	if row < 0 || col < 0 || row >= len(grid) || col >= len(grid[row]) {
		return 0
	}

	if grid[row][col] >= height {
		return 1
	}

	switch dir {
	case LEFT:
		return 1 + score(grid, height, row-1, col, dir)
	case RIGHT:
		return 1 + score(grid, height, row+1, col, dir)
	case UP:
		return 1 + score(grid, height, row, col-1, dir)
	case DOWN:
		return 1 + score(grid, height, row, col+1, dir)
	}

	panic(dir)
}

/*
func isVisible(grid [][]int, row, col int) bool {
	height := grid[row][col]

	left := true
	for c := 0; c < col; c++ {
		if grid[row][c] >= height {
			left = false
			break
		}
	}
	if left {
		return true
	}

	right := true
	for c := col + 1; c < len(grid[row]); c++ {
		if grid[row][c] >= height {
			right = false
			break
		}
	}
	if right {
		return true
	}

	up := true
	for r := 0; r < row; r++ {
		if grid[r][col] >= height {
			up = false
			break
		}
	}
	if up {
		return true
	}

	down := true
	for r := row + 1; r < len(grid); r++ {
		if grid[r][col] >= height {
			down = false
			break
		}
	}
	if down {
		return true
	}

	return false
}
*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	elves := []int{0}

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			elves = append(elves, 0)
			continue
		}

		i, err := strconv.Atoi(line)
		if err != nil {
			panic(err)
		}

		elves[len(elves)-1] += i
	}

	sum := 0
	for i := 0; i < 3; i++ {
		idx := getLargest(elves)
		sum += elves[idx]
		fmt.Println(idx, elves[idx])
		elves = append(elves[:idx], elves[idx+1:]...)
	}
	fmt.Println(sum)
}

func getLargest(a []int) int {
	curr := 0
	idx := -1
	for i, n := range a {
		if n > curr {
			curr = n
			idx = i
		}
	}

	return idx
}

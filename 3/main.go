package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	// part1(scanner)
	part2(scanner)
}

func part1(scanner *bufio.Scanner) {
	total := 0
	for scanner.Scan() {
		line := scanner.Text()
		dup := findDup(line)
		fmt.Println(string(dup), priority(dup))
		total += priority(dup)
	}

	fmt.Println(total)
}

func part2(scanner *bufio.Scanner) {
	total := 0
	idx := 0
	badges := make(map[rune]int)
	for scanner.Scan() {
		line := scanner.Text()

		seen := make(map[rune]struct{})
		for _, c := range line {
			seen[c] = struct{}{}
		}

		for k := range seen {
			badges[k] = badges[k] + 1
		}

		if idx%3 == 2 {
			for k, v := range badges {
				if v == 3 {
					fmt.Println(string(k), v)
					total += priority(k)
					break
				}
			}
			badges = map[rune]int{}
		}
		idx += 1
	}

	fmt.Println(total)
}

func priority(c rune) int {
	if c > 'Z' {
		return int(c) - int('a') + 1
	}
	return int(c) - int('A') + 27
}

func findDup(rucksack string) rune {
	first := rucksack[:len(rucksack)/2]
	second := rucksack[len(rucksack)/2:]

	if len(first) != len(second) {
		panic(len(rucksack))
	}

	seen := make(map[rune]struct{}, len(first))
	for _, c := range first {
		seen[c] = struct{}{}
	}

	for _, c := range second {
		if _, ok := seen[c]; ok {
			return c
		}
	}

	return rune(0)
}

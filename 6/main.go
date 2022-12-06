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
	for scanner.Scan() {
		line := scanner.Text()

		mark := findMarker(line, 4)
		fmt.Println(mark)
		fmt.Println(line[:mark])

		mark = findMarker(line, 14)
		fmt.Println(mark)
		fmt.Println(line[:mark])
	}
}

func findMarker(s string, sz int) int {
	for i := range s {
		if i > len(s)-sz {
			return -1
		}

		chars := s[i : i+sz]
		if allUnique(chars) {
			return i + sz
		}
	}

	return -1
}

func allUnique(s string) bool {
	seen := make(map[rune]struct{})
	for _, c := range s {
		_, ok := seen[c]
		if ok {
			return false
		}

		seen[c] = struct{}{}
	}

	return true
}

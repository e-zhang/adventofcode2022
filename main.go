package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Section struct {
	start int
	end   int
}

func (s Section) Contains(other Section) bool {
	return s.start <= other.start && s.end >= other.end
}

func (s Section) Overlaps(other Section) bool {
	return s.Within(other.start) || s.Within(other.end) || other.Within(s.start) || other.Within(s.end)
}

func (s Section) Within(idx int) bool {
	return s.start <= idx && idx <= s.end
}

func From(s string) Section {
	section := strings.Split(s, "-")
	start, err := strconv.Atoi(section[0])
	if err != nil {
		panic(err)
	}
	end, err := strconv.Atoi(section[1])
	if err != nil {
		panic(err)
	}
	return Section{
		start: start,
		end:   end,
	}
}

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
	overlaps := 0
	for scanner.Scan() {
		line := scanner.Text()
		pair := strings.Split(line, ",")
		first := From(pair[0])
		second := From(pair[1])
		if first.Contains(second) || second.Contains(first) {
			fmt.Println(pair)
			overlaps += 1
		}
	}

	fmt.Println(overlaps)
}

func part2(scanner *bufio.Scanner) {
	overlaps := 0
	for scanner.Scan() {
		line := scanner.Text()
		pair := strings.Split(line, ",")
		first := From(pair[0])
		second := From(pair[1])
		if first.Overlaps(second) {
			fmt.Println(pair)
			overlaps += 1
		}
	}

	fmt.Println(overlaps)
}

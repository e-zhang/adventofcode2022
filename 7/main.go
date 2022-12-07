package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	MAX    = 70000000
	UNUSED = 30000000
)

type File struct {
	name string
	sz   int
}

type Dir struct {
	name   string
	parent *Dir
	dirs   []*Dir
	files  []*File
}

func (d *Dir) Size() int {
	sz := 0
	for _, dir := range d.dirs {
		sz += dir.Size()
	}
	for _, f := range d.files {
		sz += f.sz
	}

	return sz
}

func PrintTree(root *Dir, prefix string) {
	fmt.Printf("%s- %s (dir) %d\n", prefix, root.name, root.Size())
	for _, d := range root.dirs {
		PrintTree(d, prefix+"  ")
	}
	for _, f := range root.files {
		fmt.Printf("%s  - %s (file, size=%d)\n", prefix, f.name, f.sz)
	}
}

func ParseTree(scanner *bufio.Scanner) *Dir {
	root := &Dir{name: "/"}
	var curr *Dir

	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, " ")

		switch tokens[0] {
		case "$":
			switch cmd := tokens[1]; cmd {
			case "cd":
				switch dir := tokens[2]; dir {
				case "/":
					curr = root
				case "..":
					curr = curr.parent
				default:
					for _, d := range curr.dirs {
						if d.name == dir {
							curr = d
							break
						}
					}

					if curr.name != dir {
						panic("cannot cd to dir that does not exist: " + dir)
					}
				}
			case "ls":
			}
		case "dir":
			curr.dirs = append(curr.dirs, &Dir{name: tokens[1], parent: curr})
		default:
			sz, err := strconv.Atoi(tokens[0])
			if err != nil {
				panic(err)
			}
			curr.files = append(curr.files, &File{name: tokens[1], sz: sz})
		}
	}
	return root
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	root := ParseTree(scanner)
	PrintTree(root, "")

	part1(root)
	part2(root)
}

func part2(root *Dir) {
	used := root.Size()
	avail := MAX - used
	fmt.Println(avail)

	var candidate *Dir
	Walk(root, func(dir *Dir) {
		if sz := dir.Size(); avail+sz > UNUSED {
			if candidate == nil || sz < candidate.Size() {
				fmt.Println(dir.name, sz)
				candidate = dir
			}
		}
	})

	fmt.Println(candidate.Size())
}

func part1(root *Dir) {
	sum := 0
	Walk(root, func(dir *Dir) {
		if sz := dir.Size(); sz < 100000 {
			fmt.Println(dir.name, dir.Size())
			sum += sz
		}
	})

	fmt.Println(sum)
}

func Walk(root *Dir, fn func(*Dir)) {
	visit := []*Dir{root}

	for len(visit) > 0 {
		next := visit[0]
		visit = visit[1:]

		fn(next)

		for _, dir := range next.dirs {
			visit = append(visit, dir)
		}

	}
}

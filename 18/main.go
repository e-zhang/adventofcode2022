package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type Coord struct {
	x int
	y int
	z int
}

func (c Coord) Add(other Coord) Coord {
	return Coord{
		c.x + other.x,
		c.y + other.y,
		c.z + other.z,
	}
}

func (c Coord) ShareFace(other Coord) bool {
	dx := abs(c.x - other.x)
	dy := abs(c.y - other.y)
	dz := abs(c.z - other.z)

	return dx+dy+dz == 1
}

type Span struct {
	min int
	max int
}

func (s Span) Contains(v int) bool {
	return s.min <= v && v <= s.max
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	var cubes []Coord
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		var x, y, z int
		_, err = fmt.Sscanf(line, "%d,%d,%d", &x, &y, &z)
		if err != nil {
			panic(err)
		}

		Coord := Coord{x, y, z}
		cubes = append(cubes, Coord)
	}

	part1(cubes)
	part2(cubes)
}

func findRange(cubes []Coord) (Span, Span, Span) {
	var x, y, z Span

	x.min, y.min, z.min = cubes[0].x, cubes[0].y, cubes[0].z
	for _, c := range cubes {
		if c.x < x.min {
			x.min = c.x
		}
		if c.x > x.max {
			x.max = c.x
		}

		if c.y < y.min {
			y.min = c.y
		}
		if c.y > y.max {
			y.max = c.y
		}

		if c.z < z.min {
			z.min = c.z
		}
		if c.z > z.max {
			z.max = c.z
		}
	}

	return Span{x.min - 1, x.max + 1}, Span{y.min - 1, y.max + 1}, Span{z.min - 1, z.max + 1}
}

func part1(cubes []Coord) {
	var shared int
	seen := map[Coord]struct{}{}
	for _, c := range cubes {
		for _, k := range cubes {
			if c == k {
				continue
			}

			if _, kok := seen[k]; kok {
				continue
			}

			if c.ShareFace(k) {
				shared++
			}
		}

		seen[c] = struct{}{}
	}

	max := 6 * len(cubes)
	surface := max - 2*shared

	fmt.Println(surface)
}

func part2(cubes []Coord) {
	x, y, z := findRange(cubes)

	lookup := make(map[Coord]struct{})
	for _, c := range cubes {
		lookup[c] = struct{}{}
	}

	q := []Coord{Coord{x.min, y.min, z.min}}
	visited := make(map[Coord]struct{})
	hits := 0
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]

		if _, ok := visited[cur]; ok {
			continue
		}

		visited[cur] = struct{}{}

		for _, delta := range []Coord{
			Coord{1, 0, 0},
			Coord{0, 1, 0},
			Coord{0, 0, 1},
			Coord{-1, 0, 0},
			Coord{0, -1, 0},
			Coord{0, 0, -1},
		} {
			next := cur.Add(delta)

			if _, ok := lookup[next]; ok {
				hits++
				continue
			}

			if x.Contains(next.x) && y.Contains(next.y) && z.Contains(next.z) {
				q = append(q, next)
			}
		}
	}

	fmt.Println(hits)
}

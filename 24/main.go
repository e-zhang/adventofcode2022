package main

import (
	"bufio"
	"fmt"
	"image"
	"os"
	"strings"
)

var debug = false

var count = make(map[string]int)

const (
	U     = '^'
	D     = 'v'
	L     = '<'
	R     = '>'
	EMPTY = '.'
	WALL  = '#'
)

type Blizzard struct {
	loc   image.Point
	delta image.Point
	dir   rune
}

func NewBlizzard(x, y int, dir rune) Blizzard {
	delta := image.Pt(0, 0)
	switch dir {
	case U:
		delta.Y = -1
	case D:
		delta.Y = 1
	case L:
		delta.X = -1
	case R:
		delta.X = 1
	default:
		panic(dir)
	}

	return Blizzard{image.Pt(x, y), delta, dir}
}

func (b *Blizzard) Move(w, h int) {
	next := b.loc.Add(b.delta)

	if next.Y >= h-1 {
		next.Y = 1
	}
	if next.Y <= 0 {
		next.Y = h - 2
	}

	if next.X >= w-1 {
		next.X = 1
	}
	if next.X <= 0 {
		next.X = w - 2
	}

	b.loc = next
}

func computeBlizzardTopology(blizzards []Blizzard, w, h int) [][]Blizzard {
	seen := make(map[string]struct{})
	combos := [][]Blizzard{blizzards}

	key := func(bs []Blizzard) string {
		var k strings.Builder
		for _, b := range bs {
			k.WriteString(b.loc.String())
		}
		return k.String()
	}

	seen[key(blizzards)] = struct{}{}

	for {
		next := make([]Blizzard, len(blizzards))
		copy(next, blizzards)
		for i := range next {
			next[i].Move(w, h)
		}
		k := key(next)
		if _, ok := seen[k]; ok {
			break
		}

		seen[k] = struct{}{}
		blizzards = next
		combos = append(combos, next)
	}

	return combos
}

type Key struct {
	pos      image.Point
	blizzard int
}

type State struct {
	start image.Point
	end   image.Point

	w int
	h int

	blizzards [][]Blizzard

	min int
}

func (s *State) IsValid(pos image.Point, time int) bool {

	if pos.Eq(s.start) || pos.Eq(s.end) {
		return true
	}

	if pos.X <= 0 || pos.X >= s.w-1 {
		return false
	}

	if pos.Y <= 0 || pos.Y >= s.h-1 {
		return false
	}

	for _, b := range s.blizzards[time%len(s.blizzards)] {
		if pos.Eq(b.loc) {
			return false
		}
	}

	return true
}

func (s *State) Print(pos image.Point, time int) {
	if !debug {
		return
	}

	fmt.Printf("=== minute %d | pos %s  ===\n", time, pos)

	for y := 0; y < s.w; y++ {
		for x := 0; x < s.h; x++ {
			pt := image.Pt(x, y)

			sq := "."
			if y == 0 || y == s.h-1 || x == 0 || x == s.w-1 {
				if !s.start.Eq(pt) && !s.end.Eq(pt) {
					sq = "#"
				}
			}

			count := 0
			for _, b := range s.blizzards[time%len(s.blizzards)] {
				if b.loc.Eq(image.Pt(x, y)) {
					sq = string(b.dir)
					count++
				}
			}
			if count > 1 {
				sq = fmt.Sprintf("%d", count)
			}

			if pos.Eq(image.Pt(x, y)) {
				sq = "E"
			}

			fmt.Printf(sq)
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	var blizzards []Blizzard

	var w, h int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		w = 0
		for _, c := range line {
			switch c {
			case U, D, L, R:
				b := NewBlizzard(w, h, c)
				blizzards = append(blizzards, b)
			case WALL, EMPTY:
				break
			}
			w++
		}
		h++
	}

	allBlizzards := computeBlizzardTopology(blizzards, w, h)
	fmt.Println(len(allBlizzards))

	start := image.Pt(1, 0)
	end := image.Pt(w-2, h-1)

	time := 0
	state := &State{
		start,
		end,
		w,
		h,
		allBlizzards,
		0,
	}
	// SimulateDFS(state, map[Key]int{}, start, time)
	SimulateBFS(state, time)
	fmt.Println(state.min)

	// part 2
	fmt.Println("=======")
	time = state.min
	state.start = end
	state.end = start
	state.min = 0
	// SimulateDFS(state, map[Key]int{}, end, time)
	SimulateBFS(state, time)
	fmt.Println("back", state.min)

	time = state.min
	state.start = start
	state.end = end
	state.min = 0
	// SimulateDFS(state, map[Key]int{}, start, time)
	SimulateBFS(state, time)
	fmt.Println("there", state.min)
}

func SimulateDFS(state *State, visited map[Key]int, pos image.Point, time int) {
	if state.min > 0 && state.min <= time {
		return
	}

	key := Key{pos, time % len(state.blizzards)}
	if t, ok := visited[key]; ok {
		if t <= time {
			return
		}
	}
	visited[key] = time

	if state.end.Eq(pos) {
		state.Print(pos, time)
		fmt.Println(state.min, time)
		state.min = time
		return
	}

	state.Print(pos, time)

	options := []image.Point{
		image.Pt(0, 1),
		image.Pt(1, 0),
		image.Pt(0, 0),
		image.Pt(-1, 0),
		image.Pt(0, -1),
	}

	if state.end.X < state.start.X {
		options = []image.Point{
			image.Pt(-1, 0),
			image.Pt(0, -1),
			image.Pt(0, 1),
			image.Pt(1, 0),
			image.Pt(0, 0),
		}
	}

	for _, opt := range options {
		next := pos.Add(opt)

		if state.IsValid(next, time+1) {
			SimulateDFS(state, visited, next, time+1)
		}
	}
}

func SimulateBFS(state *State, time int) {
	visited := make(map[Key]struct{})

	options := []image.Point{
		image.Pt(0, 1),
		image.Pt(1, 0),
		image.Pt(0, 0),
		image.Pt(-1, 0),
		image.Pt(0, -1),
	}

	q := make([]struct {
		p image.Point
		t int
	}, 0)
	q = append(q, struct {
		p image.Point
		t int
	}{state.start, time})

	for len(q) > 0 {
		cur := q[0]
		q = q[1:]

		if state.end.Eq(cur.p) {
			state.min = cur.t
			return
		}

		if _, ok := visited[Key{cur.p, cur.t % len(state.blizzards)}]; ok {
			continue
		}

		for _, opt := range options {
			next := cur.p.Add(opt)

			if state.IsValid(next, cur.t+1) {
				q = append(q, struct {
					p image.Point
					t int
				}{next, cur.t + 1})

				visited[Key{cur.p, cur.t % len(state.blizzards)}] = struct{}{}
			}
		}
	}
}

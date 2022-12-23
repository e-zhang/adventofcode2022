package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const (
	WALL  = '#'
	OPEN  = '.'
	EMPTY = ' '
)

const (
	L = 'L'
	R = 'R'
)

const (
	RIGHT = 0
	DOWN  = 1
	LEFT  = 2
	UP    = 3
)

func Dir(dir int) string {
	switch dir {
	case RIGHT:
		return "RIGHT"
	case DOWN:
		return "DOWN"
	case LEFT:
		return "LEFT"
	case UP:
		return "UP"
	}
	panic(dir)
}

type Instruction struct {
	move int
	rot  rune
}

type State struct {
	row    int
	col    int
	facing int
}

type Cell struct {
	space rune
	face  int
}

func (s *State) MovePart1(grid *Grid) bool {
	switch s.facing {
	case RIGHT, LEFT:
		step := 1
		if s.facing == LEFT {
			step = -1
		}

		c := (s.col + step + grid.Cols()) % grid.Cols()
		if grid.Get(s.row, c).space == EMPTY {
			for {
				next := (c - step + grid.Cols()) % grid.Cols()
				if grid.Get(s.row, next).space == EMPTY {
					break
				}
				c = next
			}
		}

		if grid.Get(s.row, c).space != OPEN {
			return false
		}

		s.col = c
		return true
	case DOWN, UP:
		step := 1
		if s.facing == UP {
			step = -1
		}

		r := (s.row + step + grid.Rows()) % grid.Rows()
		if grid.Get(r, s.col).space == EMPTY {
			for {
				next := (r - step + grid.Rows()) % grid.Rows()
				if grid.Get(next, s.col).space == EMPTY {
					break
				}
				r = next
			}
		}

		if grid.Get(r, s.col).space != OPEN {
			return false
		}

		s.row = r
		return true
	}

	panic(s.facing)
}

func (s *State) MovePart2(grid *Grid) bool {
	id := grid.Get(s.row, s.col).face
	face := grid.faces[id-1]

	dr := s.row - face.row
	dc := s.col - face.col

	news := *s
	switch s.facing {
	case LEFT:
		c := (dc - 1)
		news.col = face.col + c
		if c < 0 {
			next := face.neighbors[s.facing]
			facing := next.Facing(face)
			switch facing {
			case RIGHT:
				news.col = next.col + next.width - 1
				news.row = next.row + dr
			case LEFT:
				news.col = next.col
				news.row = next.row + next.height - 1 - dr
				news.facing = RIGHT
			case UP:
				news.row = next.row
				news.col = next.col + dr
				news.facing = DOWN
			case DOWN:
				news.row = next.row + next.height - 1
				news.col = next.col + next.width - 1 - dr
				news.facing = UP
			}
		}
	case RIGHT:
		c := (dc + 1)
		news.col = face.col + c
		if c >= face.width {
			next := face.neighbors[news.facing]
			facing := next.Facing(face)
			switch facing {
			case RIGHT:
				news.col = next.col + next.width - 1
				news.row = next.row + next.height - 1 - dr
				news.facing = LEFT
			case LEFT:
				news.col = next.col
				news.row = next.row + dr
			case UP:
				news.row = next.row
				news.col = next.col + next.width - 1 - dr
				news.facing = DOWN
			case DOWN:
				news.row = next.row + next.height - 1
				news.col = next.col + dr
				news.facing = UP
			}
		}
	case UP:
		r := (dr - 1)
		news.row = face.row + r
		if r < 0 {
			next := face.neighbors[news.facing]
			facing := next.Facing(face)
			switch facing {
			case LEFT:
				news.row = next.row + dc
				news.col = next.col
				news.facing = RIGHT
			case RIGHT:
				news.row = next.row + next.height - 1 - dc
				news.col = next.col + next.width - 1
				news.facing = LEFT
			case UP:
				news.row = next.row
				news.col = next.col + next.width - 1 - dc
				news.facing = DOWN
			case DOWN:
				news.row = next.row + next.height - 1
				news.col = next.col + dc
			}
		}
	case DOWN:
		r := (dr + 1)
		news.row = face.row + r
		if r >= face.height {
			next := face.neighbors[news.facing]
			facing := next.Facing(face)
			switch facing {
			case LEFT:
				news.row = next.row + next.height - 1 - dc
				news.col = next.col
				news.facing = RIGHT
			case RIGHT:
				news.row = next.row + dc
				news.col = next.col + next.width - 1
				news.facing = LEFT
			case UP:
				news.row = next.row
				news.col = next.col + dc
			case DOWN:
				news.row = next.row + next.height - 1
				news.col = next.col + next.width - 1 - dc
				news.facing = UP
			}
		}
	default:
		panic(s.facing)
	}

	if grid.Get(news.row, news.col).space != OPEN {
		return false
	}

	*s = news
	return true
}

func (s *State) Rotate(dir rune) {
	rot := 0
	switch dir {
	case L:
		rot = -1
	case R:
		rot = 1
	default:
		panic(dir)
	}
	s.facing = (s.facing + rot + 4) % 4
}

type Face struct {
	id int

	neighbors []*Face

	row    int
	col    int
	width  int
	height int
}

func (f *Face) Facing(n *Face) int {
	for i, m := range f.neighbors {
		if m == n {
			return i
		}
	}

	panic(n.id)
}

func (f *Face) Print() {
	fmt.Printf("[%d] ", f.id)

	fmt.Printf("up:")
	if f.neighbors[UP] != nil {
		fmt.Printf("%d", f.neighbors[UP].id)
	} else {
		fmt.Printf(" ")
	}
	fmt.Printf(", ")

	fmt.Printf("down:")
	if f.neighbors[DOWN] != nil {
		fmt.Printf("%d", f.neighbors[DOWN].id)
	} else {
		fmt.Printf(" ")
	}
	fmt.Printf(", ")

	fmt.Printf("left:")
	if f.neighbors[LEFT] != nil {
		fmt.Printf("%d", f.neighbors[LEFT].id)
	} else {
		fmt.Printf(" ")
	}
	fmt.Printf(", ")

	fmt.Printf("right:")
	if f.neighbors[RIGHT] != nil {
		fmt.Printf("%d", f.neighbors[RIGHT].id)
	} else {
		fmt.Printf(" ")
	}

	fmt.Println()
}

type Grid struct {
	cells [][]Cell

	faces []*Face
}

func (g *Grid) Start() *State {
	r := 0
	c := 0

	for i, s := range g.cells[r] {
		if s.space == OPEN {
			c = i
			break
		}
	}

	if c >= len(g.cells[r]) {
		panic(c)
	}

	return &State{r, c, 0}
}

func (g *Grid) Rows() int {
	return len(g.cells)
}

func (g *Grid) Cols() int {
	return len(g.cells[0])
}

func (g *Grid) Get(r, c int) Cell {
	return g.cells[r][c]
}

func (g *Grid) Print(face bool) {
	return
	for _, r := range g.cells {
		for _, c := range r {
			if face {
				if c.face == 0 {
					fmt.Printf(" ")
				} else {
					fmt.Printf("%d", c.face)
				}
			} else {
				fmt.Printf("%c", c.space)
			}
		}
		fmt.Println()
	}
}

func (g *Grid) Copy() *Grid {
	newg := Grid{}

	newg.cells = make([][]Cell, len(g.cells))
	for r, row := range g.cells {
		newg.cells[r] = make([]Cell, len(row))
		copy(newg.cells[r], row)
	}

	newg.faces = make([]*Face, len(g.faces))
	copy(newg.faces, g.faces)

	return &newg
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)

	grid := parseGrid(scanner)

	instructions := parseInstructions(scanner)

	state := grid.Start()
	grid.Print(false)
	grid.Print(true)

	grid2 := grid.Copy()
	state2 := grid2.Start()

	for _, in := range instructions {
		Walk(grid, state, in, func(s *State, g *Grid) bool {
			return s.MovePart1(g)
		})
		Walk(grid2, state2, in, func(s *State, g *Grid) bool {
			return s.MovePart2(g)
		})
	}

	fmt.Println(state)
	fmt.Println(1000*(state.row+1) + 4*(state.col+1) + state.facing)

	fmt.Println(state2)
	fmt.Println(1000*(state2.row+1) + 4*(state2.col+1) + state2.facing)
}

func Walk(grid *Grid, state *State, in Instruction, mv func(*State, *Grid) bool) {
	if in.move > 0 {
		for i := 0; i < in.move; i++ {
			if !mv(state, grid) {
				break
			}
		}
		return
	}

	if in.rot != rune(0) {
		state.Rotate(in.rot)
		return
	}

	panic(in)
}

func parseInstructions(scanner *bufio.Scanner) []Instruction {
	scanner.Scan()
	line := scanner.Text()

	var instructions []Instruction

	i := 0
	for j, c := range line {
		if c == L || c == R {
			mv, err := strconv.Atoi(line[i:j])
			if err != nil {
				panic(line[i:j])
			}

			instructions = append(instructions, Instruction{mv, rune(0)}, Instruction{0, c})
			i = j + 1
		}
	}

	if i < len(line) {
		mv, err := strconv.Atoi(line[i:])
		if err != nil {
			panic(line[i:])
		}
		instructions = append(instructions, Instruction{mv, rune(0)})
	}

	return instructions
}

func parseGrid(scanner *bufio.Scanner) *Grid {
	var grid [][]Cell
	w := 0
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			break
		}

		if len(line) > w {
			w = len(line)
		}

		row := make([]Cell, len(line))
		for i, c := range line {
			row[i].space = c
		}

		grid = append(grid, row)
	}

	for r := range grid {
		if len(grid[r]) < w {
			for i := len(grid[r]); i < w; i++ {
				grid[r] = append(grid[r], Cell{EMPTY, 0})
			}
		}
	}

	faces := labelFaces(grid)

	return &Grid{grid, faces}
}

func labelFaces(grid [][]Cell) []*Face {
	var faces []*Face

	face := 1
	w, h := faceWidthAndHeight(grid)

	for i := 0; i < len(grid)/h; i++ {
		for j := 0; j < len(grid[0])/w; j++ {
			if grid[i*h][j*w].space == EMPTY {
				continue
			}

			for r := i * h; r < (i+1)*h; r++ {
				for c := j * w; c < (j+1)*w; c++ {
					grid[r][c].face = face
				}
			}

			faces = append(faces, &Face{face, make([]*Face, 4), i * h, j * w, w, h})
			face++
		}
	}

	for _, f := range faces {
		if id := checkFaces(f, grid, [][]int{{1, 0}, {1, 1}, {1, -1}}); id > 0 {
			f.neighbors[DOWN] = faces[id-1]
		}

		if id := checkFaces(f, grid, [][]int{{-1, 0}, {-1, 1}, {-1, -1}}); id > 0 {
			f.neighbors[UP] = faces[id-1]
		}

		if id := checkFaces(f, grid, [][]int{{0, 1}, {1, 1}, {-1, 1}}); id > 0 {
			f.neighbors[RIGHT] = faces[id-1]
		}

		if id := checkFaces(f, grid, [][]int{{0, -1}, {1, -1}, {-1, -1}}); id > 0 {
			f.neighbors[LEFT] = faces[id-1]
		}
	}

	edges := 0
	for edges < len(faces)*4 {
		edges = 0
		for _, f := range faces {
			for i, n := range f.neighbors {
				if n == nil {
					continue
				}

				loopFaces(f, n, i)
				edges++
			}
		}
	}

	return faces
}

func faceWidthAndHeight(grid [][]Cell) (int, int) {
	width := len(grid[0])
	for _, r := range grid {
		start := 0
		for ; start < len(r); start++ {
			if r[start].space != EMPTY {
				break
			}
		}
		end := start + 1
		for ; end < len(r); end++ {
			if r[end].space == EMPTY {
				break
			}
		}
		w := end - start
		if w < width {
			width = w
		}
	}

	height := len(grid)
	for i := 0; i < len(grid[0]); i++ {
		start := 0
		for ; start < len(grid); start++ {
			if grid[start][i].space != EMPTY {
				break
			}
		}
		end := start + 1
		for ; end < len(grid); end++ {
			if grid[end][i].space == EMPTY {
				break
			}
		}

		h := end - start
		if h < height {
			height = h
		}
	}

	return width, height
}

func loopFaces(face, next *Face, dir int) {
	orig := face
	ret := (dir + 3) % 4

	for steps := 0; steps < 2; steps++ {
		if next == nil {
			return
		}

		r := calcFaceRotation(face, next)

		n := (4 + dir - 1 + r) % 4
		// fmt.Printf("f:%d  v:%d from:%s expected:%s dir:%s rot:%d  n:%s\n", face.id, next.id, Dir(from), Dir(expected), Dir(dir), r, Dir(n))

		dir = n
		face = next
		next = next.neighbors[dir]
	}

	face.neighbors[dir] = orig
	orig.neighbors[ret] = face
	// fmt.Println(face.id, Dir(dir), face.neighbors[dir].id)
	// fmt.Println(orig.id, Dir(ret), orig.neighbors[ret].id)
}

func checkFaces(f *Face, grid [][]Cell, delta [][]int) int {
	for _, d := range delta {
		r := f.row + f.height*d[0]
		c := f.col + f.width*d[1]

		if r < 0 || r >= len(grid) || c >= len(grid[r]) || c < 0 {
			continue
		}

		if id := grid[r][c].face; id > 0 {
			return id
		}
	}

	return -1
}

func calcFaceRotation(from *Face, to *Face) int {
	facing := from.Facing(to)
	reverse := to.Facing(from)

	expected := (2 + facing) % 4
	return (reverse - expected + 4) % 4
}

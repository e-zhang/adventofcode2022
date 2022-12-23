package main

import (
	"bufio"
	"fmt"
	"image"
	"os"
)

var debug = false

const (
	ELF   = '#'
	EMPTY = '.'
)

type Dir int

func (d Dir) String() string {
	switch d {
	case N:
		return "N"
	case E:
		return "E"
	case S:
		return "S"
	case W:
		return "W"
	}

	panic(d)
}

const (
	N Dir = iota
	E
	S
	W
)

type Cell struct {
	elf *Elf
}

type Grid struct {
	cells [][]Cell
}

func (g *Grid) Get(x, y int) Cell {
	if y >= len(g.cells) || y < 0 {
		return Cell{}
	}

	if x >= len(g.cells[y]) || x < 0 {
		return Cell{}
	}

	return g.cells[y][x]
}

func (g *Grid) Set(x, y int, e *Elf) {
	if y >= g.Rows() {
		g.resize(0, 1)
	}

	if y < 0 {
		g.resize(0, -1)
		y = 0
	}

	if x >= g.Cols() {
		g.resize(1, 0)
	}

	if x < 0 {
		g.resize(-1, 0)
		x = 0
	}

	g.cells[y][x].elf = e
}

func (g *Grid) Clear(x, y int) {
	g.cells[y][x].elf = nil
}

func (g *Grid) Rows() int {
	return len(g.cells)
}

func (g *Grid) Cols() int {
	return len(g.cells[0])
}

func (g *Grid) EmptyCount() int {
	count := 0

	for _, row := range g.cells {
		for _, cell := range row {
			if cell.elf == nil {
				count++
			}
		}
	}

	return count
}

func (g *Grid) Print() {
	if !debug {
		return
	}

	for _, row := range g.cells {
		for _, cell := range row {
			if cell.elf == nil {
				fmt.Printf(".")
			} else {
				fmt.Printf("#")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func (g *Grid) resize(x, y int) {
	if y > 0 {
		for i := 0; i < y; i++ {
			g.cells = append(g.cells, make([]Cell, g.Cols()))
		}
	}
	if y < 0 {
		for i := 0; i < -y; i++ {
			g.cells = append([][]Cell{make([]Cell, g.Cols())}, g.cells...)
		}
	}

	if x > 0 {
		for y := range g.cells {
			for i := 0; i < x; i++ {
				g.cells[y] = append(g.cells[y], Cell{})
			}
		}
	}

	if x < 0 {
		for y := range g.cells {
			for i := 0; i < -x; i++ {
				g.cells[y] = append([]Cell{Cell{}}, g.cells[y]...)
			}
		}
	}

	if x < 0 || y < 0 {
		delta := image.Pt(x, y)
		for y, row := range g.cells {
			for x, cell := range row {
				if cell.elf == nil {
					continue
				}

				p := cell.elf.pos.Sub(delta)
				if !p.Eq(image.Pt(x, y)) {
					panic(p)
				}
				cell.elf.pos = p

				if cell.elf.proposal != nil {
					p = cell.elf.proposal.Sub(delta)
					cell.elf.proposal = &p
				}
			}
		}
	}

}

func (g *Grid) HasNeighbors(e *Elf, dir Dir) bool {
	for _, off := range getOffsets(dir) {
		p := e.pos.Add(off)
		cell := g.Get(p.X, p.Y)
		if cell.elf != nil {
			return true
		}
	}
	return false
}

type Elf struct {
	pos image.Point

	proposal *image.Point
}

func (e *Elf) Move(grid *Grid) {
	if e.proposal == nil {
		return
	}

	grid.Set(e.proposal.X, e.proposal.Y, e)
	grid.Clear(e.pos.X, e.pos.Y)

	e.pos = *e.proposal
	e.proposal = nil
}

func (e *Elf) Propose(grid *Grid, dirs []Dir) bool {
	adjacent := false
	for _, d := range dirs {
		if grid.HasNeighbors(e, d) {
			adjacent = true
			break
		}
	}
	if !adjacent {
		return false
	}

	for _, d := range dirs {
		if !grid.HasNeighbors(e, d) {
			var pt image.Point
			switch d {
			case N:
				pt = e.pos.Add(image.Pt(0, -1))
			case E:
				pt = e.pos.Add(image.Pt(1, 0))
			case S:
				pt = e.pos.Add(image.Pt(0, 1))
			case W:
				pt = e.pos.Add(image.Pt(-1, 0))
			default:
				panic(d)
			}

			e.proposal = &pt
			return true
		}
	}

	return false
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	grid := &Grid{cells: [][]Cell{}}
	var elves []*Elf

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		row := make([]Cell, len(line))

		for i, c := range line {
			if c == ELF {
				e := &Elf{
					pos: image.Pt(i, grid.Rows()),
				}

				elves = append(elves, e)
				row[i].elf = e
			}
		}

		grid.cells = append(grid.cells, row)
	}

	dirs := []Dir{N, S, W, E}

	grid.Print()

	for round := 1; ; round++ {
		proposed := make(map[image.Point]*Elf)

		hasProposals := false
		for _, e := range elves {
			hasProposals = e.Propose(grid, dirs) || hasProposals

			if e.proposal != nil {
				if other, ok := proposed[*e.proposal]; ok {
					e.proposal = nil
					other.proposal = nil
				} else {
					proposed[*e.proposal] = e
				}
			}
		}

		if !hasProposals {
			fmt.Println(round)
			break
		}

		for _, e := range elves {
			e.Move(grid)
		}

		grid.Print()
		dirs = append(dirs[1:], dirs[0])
	}

	fmt.Println(grid.EmptyCount())
}

func getOffsets(dir Dir) []image.Point {
	switch dir {
	case N:
		return []image.Point{
			image.Pt(0, -1),
			image.Pt(1, -1),
			image.Pt(-1, -1),
		}
	case E:
		return []image.Point{
			image.Pt(1, 0),
			image.Pt(1, -1),
			image.Pt(1, 1),
		}
	case S:
		return []image.Point{
			image.Pt(0, 1),
			image.Pt(1, 1),
			image.Pt(-1, 1),
		}
	case W:
		return []image.Point{
			image.Pt(-1, 0),
			image.Pt(-1, -1),
			image.Pt(-1, 1),
		}
	}

	panic(dir)
}

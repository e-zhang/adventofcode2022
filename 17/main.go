package main

import (
	"bufio"
	"fmt"
	"image"
	"os"
	"strings"
)

const draw = false

const TRILLION = 1_000_000_000_000

const (
	LEFT  = '<'
	RIGHT = '>'
)

type Pattern []image.Point

var patterns = []Pattern{
	{
		image.Pt(0, 0),
		image.Pt(1, 0),
		image.Pt(2, 0),
		image.Pt(3, 0),
	}, {
		image.Pt(1, 0),
		image.Pt(0, 1),
		image.Pt(1, 1),
		image.Pt(2, 1),
		image.Pt(1, 2),
	}, {
		image.Pt(0, 0),
		image.Pt(1, 0),
		image.Pt(2, 0),
		image.Pt(2, 1),
		image.Pt(2, 2),
	}, {
		image.Pt(0, 0),
		image.Pt(0, 1),
		image.Pt(0, 2),
		image.Pt(0, 3),
	}, {

		image.Pt(0, 0),
		image.Pt(1, 0),
		image.Pt(0, 1),
		image.Pt(1, 1),
	},
}

type Rock struct {
	pattern Pattern
	pos     image.Point
}

func (r *Rock) MaxY() int {
	y := 0
	for _, pat := range r.pattern {
		if pat.Y > y {
			y = pat.Y
		}
	}
	return r.pos.Y + y
}

func (r *Rock) Fall(chamber *Chamber) bool {
	if r.pos.Y > len(chamber.units) {
		r.pos.Y--
		return true
	}

	if len(chamber.units) == 0 {
		return false
	}

	for _, pat := range r.pattern {
		pt := r.pos.Add(pat)
		if pt.Y > len(chamber.units) {
			continue
		}
		if pt.Y == 0 {
			return false
		}
		if chamber.units[pt.Y-1][pt.X] != '.' {
			return false
		}
	}

	r.pos.Y--
	return true
}

func (r *Rock) Jet(c *Chamber, j rune) {
	switch j {
	case LEFT:
		if r.pos.X == 0 {
			return
		}

		for _, pat := range r.pattern {
			pt := r.pos.Add(pat)
			if pt.Y < len(c.units) {
				if c.units[pt.Y][pt.X-1] != '.' {
					return
				}
			}
		}
		r.pos.X--

	case RIGHT:
		for _, pat := range r.pattern {
			pt := r.pos.Add(pat)
			if pt.X >= 6 {
				return
			}
			if pt.Y < len(c.units) {
				if c.units[pt.Y][pt.X+1] != '.' {
					return
				}
			}
		}
		r.pos.X++
	}
}

func Next(h int, p Pattern) *Rock {
	return &Rock{p, image.Pt(2, h+3)}
}

type Chamber struct {
	units [][]rune
}

func (c *Chamber) AddRow() {
	r := make([]rune, 7)
	for i := range r {
		r[i] = '.'
	}
	c.units = append(c.units, r)
}

func (c *Chamber) Fill(rock *Rock) {
	for _, pat := range rock.pattern {
		pt := rock.pos.Add(pat)
		for len(c.units) <= pt.Y {
			c.AddRow()
		}

		c.units[pt.Y][pt.X] = '#'
	}
}

func (c *Chamber) Tops() []int {
	tops := make([]int, 7)
	for x := range tops {
		for y := len(c.units) - 1; y >= 0; y-- {
			if c.units[y][x] != '.' {
				tops[x] = len(c.units) - y - 1
				break
			}
		}
	}

	return tops
}

func (c *Chamber) Draw(rock *Rock, n int) {
	if !draw {
		return
	}

	maxy := len(c.units) - 1
	if rock != nil && rock.MaxY() > maxy {
		maxy = rock.MaxY()
	}
	for y := maxy; y >= n; y-- {
		for x := 0; x < 7; x++ {
			char := '.'
			if y < len(c.units) {
				char = c.units[y][x]
			}
			if rock != nil {
				for _, pat := range rock.pattern {
					pt := rock.pos.Add(pat)
					if pt.X == x && pt.Y == y {
						char = '@'
					}
				}

			}
			fmt.Printf("%c", char)
		}
		fmt.Println()
	}

	fmt.Println()
}

type key struct {
	block int
	jet   int
	tops  string
}

func toString(t []int) string {
	s := make([]string, len(t))
	for i := range s {
		s[i] = fmt.Sprintf("%d", t[i])
	}
	return strings.Join(s, ",")
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	var jets string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		jets += line
	}

	chamber := &Chamber{}
	p := 0
	j := 0
	c := rune(0)

	var heights []int
	state := map[key]int{}

	var calcH int

	for i := 0; i < 2022; i++ {
		rock := Next(len(chamber.units), patterns[p])
		// chamber.Draw(rock, 0)

		for {
			c, j = nextJet(jets, j)
			rock.Jet(chamber, c)
			// chamber.Draw(rock, 0)

			if !rock.Fall(chamber) {
				break
			}
			// chamber.Draw(rock, 0)
		}

		chamber.Fill(rock)
		// chamber.Draw(rock, 0)
		p++
		if p >= len(patterns) {
			p = 0
		}

		// check for cycle
		if calcH == 0 {
			tops := chamber.Tops()
			k := key{p, j, toString(tops)}
			if v, ok := state[k]; ok {
				calcH = calcCycle(chamber, heights, v, i)
				chamber.Draw(nil, len(chamber.units)-10)
			}

			state[k] = i
			heights = append(heights, len(chamber.units))
		}
	}

	fmt.Println(len(chamber.units))
	fmt.Println(calcH)

}

func nextJet(jets string, j int) (rune, int) {
	c := rune(jets[j])
	j++
	if j >= len(jets) {
		j = 0
	}
	return c, j
}

func calcCycle(chamber *Chamber, heights []int, prev, curr int) int {
	dh := len(chamber.units) - heights[prev]
	period := curr - prev

	cycles := (TRILLION - prev) / period
	rem := (TRILLION - prev) % cycles

	add := heights[prev]
	if rem > 0 {
		// calculate the remaining pieces based on how far into the cycle they are
		add += heights[prev+rem-1] - heights[prev]
	}
	h := cycles*dh + add

	fmt.Printf("cycle detected: end=%d|start=%d|height=%d|remaining=%d\n", curr, prev, dh, add)
	return h
}

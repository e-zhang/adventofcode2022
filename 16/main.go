package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Valve struct {
	id     int
	name   string
	flow   int
	valves []*Valve
}

func (v *Valve) String() string {
	s := fmt.Sprintf("%s:%d", v.name, v.flow)

	var names []string
	for _, vv := range v.valves {
		names = append(names, vv.name)
	}

	s += " [" + strings.Join(names, ",") + "]"
	return s
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	valves := make(map[string]*Valve)
	connections := make(map[string]string)

	rex := regexp.MustCompile(`Valve ([A-Z]{2}) has flow rate=([0-9]+); tunnel[s]{0,1} lead[s]{0,1} to valve[s]{0,1} (.*)`)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		sub := rex.FindAllStringSubmatch(line, -1)

		name := sub[0][1]
		flow, err := strconv.Atoi(sub[0][2])
		if err != nil {
			panic(err)
		}

		conns := sub[0][3]

		connections[name] = conns
		valves[name] = &Valve{id: len(valves), name: name, flow: flow}
	}

	for k, v := range valves {
		conns := strings.Split(connections[k], ", ")
		v.valves = make([]*Valve, len(conns))
		for i := range conns {
			v.valves[i] = valves[conns[i]]
		}
	}

	aa := valves["AA"]

	byId := make([]*Valve, len(valves))
	for _, v := range valves {
		byId[v.id] = v
	}

	distances := buildCosts(valves)

	var names []string
	for k, v := range valves {
		if v.flow > 0 {
			names = append(names, k)
		} else {
			delete(valves, k)
		}
	}

	opened := make(map[string]struct{})
	fmt.Println(pressureReleased(aa, valves, 30, opened, distances))

	state := make(map[int]int)
	visit(aa, valves, 0, 30, 0, distances, state)

	max := 0
	for _, v := range state {
		if v > max {
			max = v
		}
	}
	fmt.Println(max)

	state = make(map[int]int)
	max = 0
	visit(aa, valves, 0, 26, 0, distances, state)
	for p1, v1 := range state {
		for p2, v2 := range state {
			if (p1 & p2) == 0 {
				if v1+v2 > max {
					max = v1 + v2
				}
			}
		}
	}
	fmt.Println(max)
}

func buildCosts(valves map[string]*Valve) [][]int {
	distances := make([][]int, len(valves))
	for i := range distances {
		distances[i] = make([]int, len(valves))
		for j := range distances[i] {
			distances[i][j] = math.MaxInt32
		}
	}

	for _, v := range valves {
		for _, u := range v.valves {
			distances[v.id][u.id] = 1
		}
		distances[v.id][v.id] = 0
	}

	for k := 0; k < len(valves); k++ {
		for i := 0; i < len(valves); i++ {
			for j := 0; j < len(valves); j++ {
				cur := distances[i][j]
				if d := distances[i][k] + distances[k][j]; cur > d {
					distances[i][j] = d
				}
			}
		}
	}

	return distances
}

func visit(src *Valve, valves map[string]*Valve, opened int, time int, released int, costs [][]int, state map[int]int) {
	if time <= 0 {
		return
	}

	if r, ok := state[opened]; !ok || released > r {
		state[opened] = released
	}

	for _, v := range valves {
		if opened&(1<<v.id) != 0 {
			continue
		}

		cost := (time - costs[src.id][v.id] - 1)
		score := v.flow * cost
		visit(v, valves, opened|(1<<v.id), cost, released+score, costs, state)
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/ernestosuarez/itertools"
)

type Valve struct {
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
		valves[name] = &Valve{name: name, flow: flow}
	}

	for k, v := range valves {
		conns := strings.Split(connections[k], ", ")
		v.valves = make([]*Valve, len(conns))
		for i := range conns {
			v.valves[i] = valves[conns[i]]
		}
		fmt.Println(v.String())
	}

	aa := valves["AA"]
	costs := make(map[string]map[string]int)
	for k, v := range valves {
		for k2, v2 := range valves {
			cost := getCost(v, v2)
			if _, ok := costs[k]; !ok {
				costs[k] = make(map[string]int)
			}
			costs[k][k2] = cost
		}
	}

	opened := make(map[string]struct{})
	fmt.Println(pressureReleased(aa, valves, 30, opened, costs))

	var names []string
	for k, v := range valves {
		if v.flow > 0 {
			names = append(names, k)
		}
	}

	max := 0
	var mut sync.Mutex
	for i := 1; i < len(names); i++ {
		for v := range itertools.CombinationsStr(names, i) {
			go func() {
				mine := make(map[string]*Valve)
				elephant := make(map[string]*Valve)
				opened := make(map[string]struct{})
				for _, n := range v {
					mine[n] = valves[n]
				}
				for _, n := range names {
					if _, ok := mine[n]; !ok {
						elephant[n] = valves[n]
					}
				}
				p := pressureReleased(aa, elephant, 26, opened, costs) + pressureReleased(aa, mine, 26, opened, costs)
				mut.Lock()
				if p > max {
					max = p
				}
				mut.Unlock()
			}()
		}
	}

	fmt.Println(max)
}

func pressureReleased(source *Valve, valves map[string]*Valve, time int, opened map[string]struct{}, costs map[string]map[string]int) int {
	if time <= 0 {
		return 0
	}

	max := 0
	for _, mv := range moves(source, valves, time, opened, costs) {
		opened[mv.next.name] = struct{}{}
		p := mv.score + pressureReleased(mv.next, valves, time-mv.cost-1, opened, costs)
		delete(opened, mv.next.name)

		if p > max {
			max = p
		}
	}

	return max
}

type Move struct {
	cost  int
	score int
	next  *Valve
}

func moves(source *Valve, valves map[string]*Valve, time int, opened map[string]struct{}, costs map[string]map[string]int) []Move {
	var moves []Move
	for _, v := range valves {
		_, ok := opened[v.name]
		if v == source || ok || v.flow == 0 {
			continue
		}

		c := costs[source.name][v.name]
		score := v.flow * (time - c - 1)

		moves = append(moves, Move{c, score, v})
	}

	sort.Slice(moves, func(i, j int) bool {
		return moves[i].score < moves[j].score
	})

	return moves
}

func getCost(source *Valve, dest *Valve) int {
	q := []*Valve{source}
	distances := make(map[string]int)
	distances[source.name] = 0
	for len(q) != 0 {
		cur := q[0]
		q = q[1:]
		if cur == dest {
			break
		}

		for _, v := range cur.valves {
			_, ok := distances[v.name]
			if !ok || distances[v.name] > distances[cur.name]+1 {
				distances[v.name] = distances[cur.name] + 1
				q = append(q, v)
			}
		}
	}

	return distances[dest.name]
}

package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	inputFmt = `Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.`
)

const (
	ORE = iota
	CLAY
	OBSIDIAN
	GEODE
)

type Cost struct {
	robot int

	ores      int
	clays     int
	obsidians int
}

func (c Cost) CanBuild(ores, clays, obsidians int) bool {
	return ores >= c.ores && clays >= c.clays && obsidians >= c.obsidians
}

type State struct {
	time int

	ores      int
	clays     int
	obsidians int
	geodes    int

	oreRobots      int
	clayRobots     int
	obsidianRobots int
	geodeRobots    int

	skipped []int
}

func (s *State) ToKey() string {
	return fmt.Sprintf("%v", s)
}

func (s *State) SubtractCost(c Cost) {
	s.ores -= c.ores
	s.clays -= c.clays
	s.obsidians -= c.obsidians
}

func (s *State) Produce() {
	s.ores += s.oreRobots
	s.clays += s.clayRobots
	s.obsidians += s.obsidianRobots
	s.geodes += s.geodeRobots

	s.time--
}

func (s *State) BuildRobot(c Cost) {
	switch c.robot {
	case ORE:
		s.oreRobots++
	case CLAY:
		s.clayRobots++
	case OBSIDIAN:
		s.obsidianRobots++
	case GEODE:
		s.geodeRobots++
	}
}

func (s *State) ShouldBuild(c Cost, bp Blueprint) bool {
	skip := false
	for _, r := range s.skipped {
		if r == c.robot {
			skip = true
			break
		}
	}

	if skip {
		return false
	}

	switch c.robot {
	case ORE:
		return s.oreRobots < bp.ores.ores ||
			s.oreRobots < bp.clays.ores ||
			s.oreRobots < bp.obsidians.ores ||
			s.oreRobots < bp.geodes.ores
	case CLAY:
		return s.clayRobots < bp.obsidians.clays
	case OBSIDIAN:
		return s.obsidianRobots < bp.geodes.obsidians
	case GEODE:
		return true
	}

	panic(c)
}

type Blueprint struct {
	id int

	ores      Cost
	clays     Cost
	obsidians Cost
	geodes    Cost
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	var blueprints []Blueprint

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		var id, ore, clayOre, obsidianOre, obsidianClay, geodeOre, geodeObsidian int
		fmt.Sscanf(line, inputFmt, &id, &ore, &clayOre, &obsidianOre, &obsidianClay, &geodeOre, &geodeObsidian)

		bp := Blueprint{
			id:        id,
			ores:      Cost{ORE, ore, 0, 0},
			clays:     Cost{CLAY, clayOre, 0, 0},
			obsidians: Cost{OBSIDIAN, obsidianOre, obsidianClay, 0},
			geodes:    Cost{GEODE, geodeOre, 0, geodeObsidian},
		}

		blueprints = append(blueprints, bp)
	}

	quality := 0
	for _, bp := range blueprints {
		init := State{time: 24, oreRobots: 1}
		dp := make(map[string]int)
		geodes := simulate(bp, init, dp)
		fmt.Println(bp, geodes)
		quality += bp.id * geodes
	}

	fmt.Println(quality)

	product := 1
	for i := 0; i < 3; i++ {
		init := State{time: 32, oreRobots: 1}
		dp := make(map[string]int)
		geodes := simulate(blueprints[i], init, dp)
		fmt.Println(blueprints[i], geodes)

		product *= geodes

		if i > len(blueprints) {
			break
		}
	}

	fmt.Println(product)
}

func simulate(bp Blueprint, state State, dp map[string]int) int {
	if state.time == 0 {
		return state.geodes
	}

	if s, ok := dp[state.ToKey()]; ok {
		return s
	}

	var alts []struct {
		State
		Cost
	}

	start := state

	for _, c := range []Cost{
		bp.geodes,
		bp.obsidians,
		bp.clays,
		bp.ores,
	} {
		if !state.ShouldBuild(c, bp) {
			continue
		}

		if c.CanBuild(state.ores, state.clays, state.obsidians) {
			alts = append(alts, struct {
				State
				Cost
			}{start, c})
		}
	}

	var max int
	state.Produce()
	for _, alt := range alts {
		state.skipped = append(state.skipped, alt.Cost.robot)
	}
	max = simulate(bp, state, dp)

	for _, alt := range alts {
		alt.State.SubtractCost(alt.Cost)
		alt.State.Produce()
		alt.State.BuildRobot(alt.Cost)
		score := simulate(bp, alt.State, dp)
		if score > max {
			max = score
		}
	}

	dp[start.ToKey()] = max
	return max
}

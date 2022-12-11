package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Test struct {
	divisor    int
	trueThrow  int
	falseThrow int
}

func (t *Test) Do(item *Item) int {
	if item.IsDivisibleBy(t.divisor) {
		return t.trueThrow
	}
	return t.falseThrow
}

type Operation struct {
	op    string
	value int
}

func (o *Operation) Do(item *Item) {
	switch o.op {
	case "*":
		item.Mul(o.value)
	case "+":
		item.Add(o.value)
	case "squared":
		item.Squared()
	}
}

type Monkey struct {
	idx   int
	items []*Item

	operation Operation
	test      Test

	inspects int
}

func (m *Monkey) DoTurn(monkeys []*Monkey, relief int) {
	for m.Next() {
		item := m.Pop()
		m.inspects++
		m.operation.Do(item)
		item.Relief(relief)
		to := m.test.Do(item)
		monkeys[to].Push(item)
	}
}

func (m *Monkey) Next() bool {
	return len(m.items) > 0
}

func (m *Monkey) Pop() *Item {
	item := m.items[0]
	m.items = m.items[1:]
	return item
}

func (m *Monkey) Push(item *Item) {
	m.items = append(m.items, item)
}

func (m *Monkey) String() string {
	elems := make([]string, len(m.items))
	for i := range m.items {
		elems[i] = fmt.Sprintf("%d", m.items[i].worry)
	}
	return fmt.Sprintf("Monkey %d: %s", m.idx, strings.Join(elems, ","))
}

type Item struct {
	worry int
}

func NewItem(worry int) *Item {
	return &Item{
		worry: worry,
	}
}

func (i *Item) IsDivisibleBy(divisor int) bool {
	return i.worry%divisor == 0
}

func (i *Item) Mul(value int) {
	i.worry *= value
}

func (i *Item) Add(value int) {
	i.worry += value
}

func (i *Item) Squared() {
	i.worry *= i.worry
}

func (i *Item) Relief(factor int) {
	// part 1
	// i.worry /= 3

	i.worry %= factor
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	var monkeys []*Monkey
	var curr *Monkey

	reliefFactor := 1

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			monkeys = append(monkeys, curr)
			reliefFactor *= curr.test.divisor
			curr = nil
		} else {
			curr = Parse(curr, line)
			curr.idx = len(monkeys)
		}
	}
	monkeys = append(monkeys, curr)
	reliefFactor *= curr.test.divisor

	// part 1
	// for round := 0; round < 20; round++ {
	for round := 0; round < 10000; round++ {
		for _, monkey := range monkeys {
			monkey.DoTurn(monkeys, reliefFactor)
		}

		if (round+1)%1000 == 0 {
			fmt.Printf("=== Round %d ===\n", round+1)
			for _, monkey := range monkeys {
				fmt.Printf("Monkey %d inspected items %d times\n", monkey.idx, monkey.inspects)
			}
		}
	}

	sort.Slice(monkeys, func(i, j int) bool {
		return monkeys[j].inspects < monkeys[i].inspects
	})
	for _, monkey := range monkeys {
		fmt.Printf("Monkey %d inspected items %d times\n", monkey.idx, monkey.inspects)
	}
	fmt.Printf("monkey business: %d\n", monkeys[0].inspects*monkeys[1].inspects)
}

func Parse(m *Monkey, line string) *Monkey {
	line = strings.TrimSpace(line)
	tokens := strings.Split(line, ":")
	switch tokens[0] {
	case "Starting items":
		items := strings.Split(tokens[1], ",")
		m.items = make([]*Item, len(items))
		for i := range items {
			v, err := strconv.Atoi(strings.TrimSpace(items[i]))
			if err != nil {
				panic(err)
			}
			m.items[i] = NewItem(v)
		}
	case "Operation":
		var op, val string
		_, err := fmt.Sscanf(tokens[1], " new = old %s %s", &op, &val)
		if err != nil {
			panic(err)
		}

		var v int
		if val == "old" && op == "*" {
			op = "squared"
		} else {
			v, err = strconv.Atoi(val)
			if err != nil {
				panic(err)
			}
		}

		m.operation = Operation{op, v}
	case "Test":
		var d int
		_, err := fmt.Sscanf(tokens[1], " divisible by %d", &d)
		if err != nil {
			panic(err)
		}
		m.test.divisor = d
	case "If true":
		var idx int
		_, err := fmt.Sscanf(tokens[1], " throw to monkey %d", &idx)
		if err != nil {
			panic(err)
		}
		m.test.trueThrow = idx
	case "If false":
		var idx int
		_, err := fmt.Sscanf(tokens[1], " throw to monkey %d", &idx)
		if err != nil {
			panic(err)
		}
		m.test.falseThrow = idx
	default:
		if strings.HasPrefix(tokens[0], "Monkey") {
			if m != nil {
				panic(line)
			}
			m = &Monkey{}
		}
	}

	return m
}

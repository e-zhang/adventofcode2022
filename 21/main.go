package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	ADD = "+"
	SUB = "-"
	MUL = "*"
	DIV = "/"
	EQ  = "="
)

const (
	HUMAN = "humn"
)

type Op struct {
	operation string
	lhs       string
	rhs       string
}

type Monkey struct {
	id string

	value int
	op    *Op
}

func (m *Monkey) Resolve(monkeys map[string]Monkey, human bool) (bool, int) {
	if m.op == nil {
		return human && m.id == HUMAN, m.value
	}

	m1 := monkeys[m.op.lhs]
	xh, x := m1.Resolve(monkeys, human)
	if human && xh {
		return true, 0
	}

	m2 := monkeys[m.op.rhs]
	yh, y := m2.Resolve(monkeys, human)
	if human && yh {
		return true, 0
	}

	switch m.op.operation {
	case ADD:
		return false, x + y
	case SUB:
		return false, x - y
	case MUL:
		return false, x * y
	case DIV:
		return false, x / y
	}

	panic(m.op)
}

func (m *Monkey) Solve(monkeys map[string]Monkey, res int) int {
	if m.id == HUMAN {
		return res
	}

	m1 := monkeys[m.op.lhs]
	xh, x := m1.Resolve(monkeys, true)

	m2 := monkeys[m.op.rhs]
	yh, y := m2.Resolve(monkeys, true)

	if xh {
		switch m.op.operation {
		case ADD:
			res -= y
		case SUB:
			res += y
		case MUL:
			res /= y
		case DIV:
			res *= y
		case EQ:
			res = y
		}

		return m1.Solve(monkeys, res)
	}

	if yh {
		switch m.op.operation {
		case ADD:
			res -= x
		case SUB:
			res = x - res
		case MUL:
			res /= x
		case DIV:
			res = x / res
		case EQ:
			res = x
		}

		return m2.Solve(monkeys, res)
	}

	panic(m)
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	monkeys := make(map[string]Monkey)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		tokens := strings.Split(line, ":")
		monkey := tokens[0]

		tokens = strings.Split(strings.TrimSpace(tokens[1]), " ")
		switch len(tokens) {
		case 1:
			v, err := strconv.Atoi(tokens[0])
			if err != nil {
				panic(err)
			}

			monkeys[monkey] = Monkey{monkey, v, nil}

		case 3:
			monkeys[monkey] = Monkey{monkey, 0, &Op{tokens[1], tokens[0], tokens[2]}}

		default:
			panic(tokens)
		}
	}

	root := monkeys["root"]
	_, res := root.Resolve(monkeys, false)
	fmt.Println(res)

	root.op.operation = EQ
	fmt.Println(root.Solve(monkeys, 0))
}

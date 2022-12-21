package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const debug = false

const DECRYPTION_KEY = 811589153

type Node struct {
	v int

	next *Node
	prev *Node
}

func (n *Node) Remove() {
	n.prev.next = n.next
	n.next.prev = n.prev
}

func (n *Node) InsertLeft(other *Node) {
	other.next = n
	other.prev = n.prev

	n.prev.next = other
	n.prev = other
}

func (n *Node) InsertRight(other *Node) {
	other.prev = n
	other.next = n.next

	n.next.prev = other
	n.next = other
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	var nodes []*Node
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		num, err := strconv.Atoi(line)
		if err != nil {
			panic(err)
		}

		node := &Node{num, nil, nil}

		nodes = append(nodes, node)
	}

	// part1
	reset(nodes)
	head := nodes[0]

	head = Mix(nodes, head)
	CalcGrove(nodes)

	fmt.Println("===================")

	// part2
	reset(nodes)
	head = nodes[0]
	for _, n := range nodes {
		n.v *= DECRYPTION_KEY
	}

	Print(head)
	for i := 0; i < 10; i++ {
		head = Mix(nodes, head)
		Print(head)
	}
	CalcGrove(nodes)
}

func Print(head *Node) {
	if !debug {
		return
	}

	fmt.Printf("%d", head.v)
	cur := head.next
	for cur != head {
		fmt.Printf(", %d", cur.v)
		cur = cur.next
	}

	fmt.Println()
}

func Mix(nodes []*Node, head *Node) *Node {
	for _, node := range nodes {
		if node.v == 0 {
			continue
		}

		if node == head {
			head = node.next
		}

		steps := abs(node.v) % (len(nodes) - 1)
		dir := sign(node.v)

		node.Remove()
		cur := move(node, dir, steps)

		if dir < 0 {
			cur.InsertLeft(node)
		} else {
			cur.InsertRight(node)
		}
		// Print(head)
	}

	return head
}

func CalcGrove(nodes []*Node) {
	zero := indexOf(nodes, 0)

	coord := 0
	for i := 1; i <= 3; i++ {
		steps := (i * 1000) % len(nodes)

		cur := move(nodes[zero], 1, steps)

		coord += cur.v
		fmt.Println(i*1000, steps, cur.v)
	}

	fmt.Println(coord)
}

func indexOf(l []*Node, x int) int {
	for i := range l {
		if l[i].v == x {
			return i
		}
	}

	panic(x)
}

func sign(x int) int {
	if x < 0 {
		return -1
	}

	if x > 0 {
		return 1
	}

	return 0
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func move(node *Node, dir, n int) *Node {
	cur := node
	for j := 0; j < abs(n); j++ {
		if dir > 0 {
			cur = cur.next
		} else {
			cur = cur.prev
		}
	}
	return cur
}

func reset(nodes []*Node) {
	for i, n := range nodes {
		prev := (i - 1 + len(nodes)) % len(nodes)
		next := (i + 1) % len(nodes)

		n.prev = nodes[prev]
		n.next = nodes[next]
	}
	head := nodes[0]
	Print(head)
}

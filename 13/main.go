package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
)

const (
	IN_ORDER = iota
	NOT_IN_ORDER
	NEXT
)

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)

	var idx, count int

	var packets [][]interface{}

	for scanner.Scan() {
		var left, right []interface{}
		json.Unmarshal(scanner.Bytes(), &left)
		scanner.Scan()
		json.Unmarshal(scanner.Bytes(), &right)
		scanner.Scan()

		idx++
		cmp := compare(left, right)
		if cmp == IN_ORDER {
			count += idx
		}

		packets = append(packets, left, right)
	}

	fmt.Println(count)

	var divider2, divider6 []interface{}
	json.Unmarshal([]byte("[[2]]"), &divider2)
	json.Unmarshal([]byte("[[6]]"), &divider6)

	packets = append(packets, divider2, divider6)
	sort.Slice(packets, func(i, j int) bool {
		return compare(packets[i], packets[j]) == IN_ORDER
	})

	distress := 1
	for i, p := range packets {
		if reflect.DeepEqual(p, divider2) || reflect.DeepEqual(p, divider6) {
			distress *= (i + 1)
		}
	}

	fmt.Println(distress)
}

func compare(left, right []interface{}) int {
	for i := range left {
		if i >= len(right) {
			return NOT_IN_ORDER
		}

		cmp := compareValue(left[i], right[i])
		// fmt.Printf("Compare %v vs %v: %v\n", left[i], right[i], cmp)
		if cmp != NEXT {
			return cmp
		}
	}

	if len(left) < len(right) {
		return IN_ORDER
	} else if len(left) > len(right) {
		return NOT_IN_ORDER
	}
	return NEXT
}

func compareValue(left, right interface{}) int {
	lIntVal, lIntOK := left.(float64)
	rIntVal, rIntOK := right.(float64)

	lListVal, lListOK := left.([]interface{})
	rListVal, rListOK := right.([]interface{})

	switch {
	case lIntOK && rIntOK:
		if lIntVal < rIntVal {
			return IN_ORDER
		} else if lIntVal > rIntVal {
			return NOT_IN_ORDER
		} else {
			return NEXT
		}
	case lListOK && rListOK:
		return compare(lListVal, rListVal)
	case lIntOK && rListOK:
		return compare([]interface{}{lIntVal}, rListVal)
	case lListOK && rIntOK:
		return compare(lListVal, []interface{}{rIntVal})
	}

	panic(fmt.Sprintf("%#v - %#v", left, right))
}

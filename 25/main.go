package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	TWO          = '2'
	ONE          = '1'
	ZERO         = '0'
	MINUS        = '-'
	DOUBLE_MINUS = '='
)

const BASE = 5

type SNAFU string

func (s SNAFU) Decimal() int {
	num := 0

	pow := 1
	for i := len(s) - 1; i >= 0; i-- {
		var digit int

		switch s[i] {
		case TWO:
			digit = 2
		case ONE:
			digit = 1
		case ZERO:
			digit = 0
		case MINUS:
			digit = -1
		case DOUBLE_MINUS:
			digit = -2
		default:
			panic(s)
		}

		num += digit * pow
		pow *= BASE
	}

	return num
}

func From(n int) SNAFU {
	s := ""

	for n > 0 {
		digit := rune(0)
		switch n % BASE {
		case 0:
			digit = ZERO
		case 1:
			digit = ONE
		case 2:
			digit = TWO
		case 3:
			digit = DOUBLE_MINUS
			n += BASE
		case 4:
			digit = MINUS
			n += BASE
		default:
			panic(n)
		}

		s = string(digit) + s
		n /= BASE
	}

	return SNAFU(s)
}

func main() {
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	sum := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		s := SNAFU(line)
		n := s.Decimal()
		fmt.Println(line, n)
		sum += n
	}

	fmt.Println(sum)

	s := From(sum)
	fmt.Println(s)

	if s.Decimal() != sum {
		panic(s)
	}
}

package main

import (
	"bufio"
	"fmt"
	"image"
	"os"
	"sort"
)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type Span struct {
	start int
	end   int
}

func distance(a, b image.Point) int {
	diff := a.Sub(b)
	return abs(diff.X) + abs(diff.Y)
}

func main() {
	// f, err := os.Open("test")
	f, err := os.Open("input")
	if err != nil {
		panic(err)
	}

	var sensors, beacons []image.Point

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		var sx, sy, bx, by int
		fmt.Sscanf(line, "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &sx, &sy, &bx, &by)

		sensor := image.Pt(sx, sy)
		sensors = append(sensors, sensor)
		beacon := image.Pt(bx, by)
		beacons = append(beacons, beacon)
	}

	// part 1
	spans := findSpans(sensors, beacons, 2000000)
	part1(spans)

	for y := 0; y < 4000000; y++ {
		// for y := 0; y < 20; y++ {
		spans := findSpans(sensors, beacons, y)
		x := 0
		for _, sp := range spans {
			if sp.start <= x && x <= sp.end {
				x = sp.end + 1
			}

			if x < sp.start {
				break
			}
		}

		if x <= 4000000 {
			fmt.Println(x, y)
			fmt.Println(x*4000000 + y)
		}

		// if !contained {
		// 	fmt.Println("spot", x, y)
		// 	fmt.Println(x*4000000 + y)
		// }
	}
}

func findSpans(sensors, beacons []image.Point, y int) []Span {
	spans := []Span{}
	for i := range sensors {
		s := sensors[i]
		b := beacons[i]
		dist := distance(s, b)

		dy := abs(s.Y - y)
		dx := dist - dy

		if dx >= 0 {
			span := Span{s.X - dx, s.X + dx}
			spans = append(spans, span)
		}
	}
	sort.Slice(spans, func(i, j int) bool {
		return spans[i].start < spans[j].start
	})
	return spans
}

func part1(spans []Span) {
	var span Span
	for _, s := range spans {
		if span.start > s.end || span.end < s.start {
			fmt.Println("skipped", span, s)
			continue
		}
		span = Span{min(span.start, s.start), max(span.end, s.end)}
	}

	fmt.Println(span.end - span.start)
}

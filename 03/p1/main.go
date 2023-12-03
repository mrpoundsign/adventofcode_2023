package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	STATE_SEARCH = iota
	STATE_NUMBER
)

type point struct {
	x, y int
}

type area struct {
	topLeft, bottomRight point
}

func (a area) intersects(b area) bool {
	if a.topLeft.x > b.bottomRight.x || a.bottomRight.x < b.topLeft.x {
		return false
	}
	if a.topLeft.y > b.bottomRight.y || a.bottomRight.y < b.topLeft.y {
		return false
	}
	return true
}

func (a area) String() string {
	return fmt.Sprintf("(%d,%d)-(%d,%d)", a.topLeft.x, a.topLeft.y, a.bottomRight.x, a.bottomRight.y)
}

type symbol struct {
	area
	symbol byte
}

func (s symbol) String() string {
	return fmt.Sprintf("%c:%s", s.symbol, s.area.String())
}

type part struct {
	area
	number int
}

func (p part) String() string {
	return fmt.Sprintf("%d:%s", p.number, p.area.String())
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal("must specify one file")
	}

	file := flag.Arg(0)
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	parts := []part{}
	symbols := []symbol{}
	scanner := bufio.NewScanner(f)
	s := STATE_SEARCH
	lineN := -1
	curr, ns, ne, nl := 0, 0, 0, 0

	for scanner.Scan() {
		lineN++
		s = STATE_SEARCH

		line := scanner.Bytes()
		for i, c := range line {
		process:
			switch s {
			case STATE_SEARCH:
				if curr != 0 {
					parts = append(parts, part{
						number: curr,
						area: area{
							topLeft: point{
								x: ns,
								y: nl,
							},
							bottomRight: point{
								x: ne,
								y: nl,
							},
						},
					})
					curr = 0
				}

				if c == '.' {
					continue
				}

				if c >= '0' && c <= '9' {
					curr = int(c - '0')
					ns = i
					ne = i
					nl = lineN
					s = STATE_NUMBER
					continue
				}

				symbols = append(symbols, symbol{
					symbol: c,
					area: area{
						topLeft: point{
							x: i - 1,
							y: lineN - 1,
						},
						bottomRight: point{
							x: i + 1,
							y: lineN + 1,
						},
					},
				})
			case STATE_NUMBER:
				if c < '0' || c > '9' {
					s = STATE_SEARCH
					goto process
				}

				ne++
				curr *= 10
				curr += int(c - '0')
			}
		}
	}

	log.Printf("symbols: %v", symbols)
	log.Printf("parts: %v", parts)

	valid_parts := []int{}
	valid_parts_sum := 0
	for _, b := range parts {
		for _, a := range symbols {
			if !a.intersects(b.area) {
				continue
			}
			valid_parts = append(valid_parts, b.number)
			valid_parts_sum += b.number
			break
		}
	}

	log.Printf("valid parts: %d, %v", valid_parts_sum, valid_parts)
}

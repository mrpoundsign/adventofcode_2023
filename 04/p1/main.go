package main

import (
	"bufio"
	"flag"
	"log"
	"math"
	"os"
)

const (
	STATE_BEGIN = iota
	STATE_NUMBERS
	STATE_WINNERS
)

type card struct {
	id      int
	numbers []int
	winners []int
}

func (c card) isWinMatch(n int) bool {
	for _, i := range c.winners {
		if i == n {
			return true
		}
	}
	return false
}

func (c card) points() int {
	var score float64 = 0

	for _, i := range c.numbers {
		if !c.isWinMatch(i) {
			continue
		}

		score += 1
	}

	if score < 2 {
		return int(score)
	}

	// log.Printf("+card %d has %d, %d points", c.id, int(score), int(math.Pow(2, score-1)))

	return int(math.Pow(2, score-1))
}

func cardFromLine(l []byte) *card {
	c := card{
		numbers: []int{},
		winners: []int{},
	}

	s := STATE_BEGIN
	cn := 0
	dc := 0

	for _, ch := range l {
		// fill current number if we're reading it
		if ch <= '9' && ch >= '0' {
			dc++
			cn *= 10
			cn += int(ch - '0')
			continue
		}

		switch s {
		case STATE_BEGIN:
			if ch == ':' {
				c.id = cn
				cn = 0
				s = STATE_NUMBERS
				continue
			}
		case STATE_NUMBERS:
			if ch != ' ' && ch != '|' {
				continue
			}

			c.numbers = append(c.numbers, cn)
			cn = 0

			if ch == '|' {
				s = STATE_WINNERS
			}
		case STATE_WINNERS:
			if ch != ' ' || cn == 0 {
				continue
			}

			c.winners = append(c.winners, cn)
			cn = 0
		}
	}
	if cn > 0 {
		c.winners = append(c.winners, cn)
	}
	return &c
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

	points := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l := scanner.Bytes()
		c := cardFromLine(l)

		points += c.points()
	}

	log.Println(points)
}

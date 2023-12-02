package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"strconv"
)

const (
	STATE_GAME       int = iota // 0
	STATE_FIND_ROLL             // 1
	STATE_BLUE_ROLL             // 2
	STATE_RED_ROLL              // 3
	STATE_GREEN_ROLL            // 4
	STATE_END                   // 5
)

type game struct {
	blue, red, green int
}

func main() {
	var games = map[int][]game{}

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

	scanner := bufio.NewScanner(f)
	currIdx := -1
	curr := game{}
	buff := []byte{}

	for scanner.Scan() {
		s := STATE_GAME
		for _, c := range scanner.Bytes() {
			addCharDigits(&buff, c)

			switch s {
			case STATE_GAME:
				switch c {
				case ':':
					i, err := strconv.Atoi(string(buff))
					if err != nil {
						log.Fatal(err)
					}
					buff = []byte{}
					currIdx = i
					fallthrough
				case ';':
					s = STATE_FIND_ROLL
				}
			case STATE_FIND_ROLL:
				switch c {
				case 'b':
					s = STATE_BLUE_ROLL
				case 'r':
					s = STATE_RED_ROLL
				case 'g':
					s = STATE_GREEN_ROLL
				}
			default:
				if c != ',' && c != ';' {
					continue
				}
				addRoll(&curr, s, &buff)
				games[currIdx] = append(games[currIdx], curr)
				curr = game{}
				s = STATE_FIND_ROLL
			}
		}

		addRoll(&curr, s, &buff)
		games[currIdx] = append(games[currIdx], curr)
		curr = game{}
	}

	powers := 0
	for _, g := range games {
		var rr, rg, rb int = 0, 0, 0
		for _, t := range g {
			if rr < t.red {
				rr = t.red
			}
			if rg < t.green {
				rg = t.green
			}
			if rb < t.blue {
				rb = t.blue
			}
		}
		powers += rr * rg * rb
	}
	log.Println(powers)
}

func addCharDigits(buff *[]byte, c byte) {
	if c < '0' || c > '9' {
		return
	}
	*buff = append(*buff, c)
}

func addRoll(g *game, s int, buff *[]byte) {
	i, err := strconv.Atoi(string(*buff))
	if err != nil {
		log.Fatal(err)
	}

	switch s {
	case STATE_BLUE_ROLL:
		g.blue += i
	case STATE_RED_ROLL:
		g.red += i
	case STATE_GREEN_ROLL:
		g.green += i
	default:
		return
	}

	*buff = []byte{}
}

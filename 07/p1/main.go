package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/mrpoundsign/adventofcode_2023/07/game"
)

const ()

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

	hands := game.Hands{}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		hands = append(hands, *game.NewHand(fields[0], fields[1]))
	}

	sort.Sort(hands)

	winnings := 0
	for i, h := range hands {
		winnings += h.Bid() * (i + 1)
	}

	log.Println(winnings)
}

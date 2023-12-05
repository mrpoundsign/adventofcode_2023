package main

import (
	"flag"
	"log"
	"os"

	"github.com/mrpoundsign/adventofcode_2023/05/almanac"
)

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

	parser := almanac.NewParser(f)
	seedList, seedMaps := parser.Parse()

	lowest := 0
	for _, s := range seedList.Seeds {
		from := "seed"
		si := s
		for _, m := range seedMaps {
			if m.From != from {
				log.Println(m.From, "no from", from)
				continue
			}

			from = m.To
			si = m.LocationFor(si)
		}
		if lowest == 0 || lowest > si {
			lowest = si
		}
	}
	log.Printf("lowest %d", lowest)
}

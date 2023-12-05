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

	// I know this is slow. Might work on it later.
	for i := 0; i < len(seedList.Seeds); i += 2 {
		for j := 0; j < seedList.Seeds[i+1]; j++ {
			s := seedList.Seeds[i] + j
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
			if lowest > si || lowest == 0 {
				lowest = si
			}
		}
	}
	log.Printf("lowest %d", lowest)
}

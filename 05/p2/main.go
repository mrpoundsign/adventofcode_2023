package main

import (
	"flag"
	"log"
	"math"
	"os"
	"sync"

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

	// Checking order of maps
	from := "seed"
	for _, m := range seedMaps {
		if m.From != from {
			log.Println(m.From, "no from", from)
		}
		from = m.To
	}

	lowest := math.MaxInt

	lowChan := make(chan int)
	var wg sync.WaitGroup

	// I know this is slow. Might work on it later.
	for i := 0; i < len(seedList.Seeds); i += 2 {
		wg.Add(1)
		go func(i, jl int) {
			defer wg.Done()

			lowest := math.MaxInt

			for j := 0; j < jl; j++ {

				s := seedList.Seeds[i] + j
				si := s

				for mi := range seedMaps {
					si = seedMaps[mi].LocationFor(si)
				}

				if lowest > si {
					lowest = si
				}
			}

			lowChan <- lowest
		}(i, seedList.Seeds[i+1])
	}

	wait := make(chan struct{})
	go func() {
		wg.Wait()
		close(wait)
	}()

	for {
		select {
		case <-wait:
			close(lowChan)
			log.Printf("lowest %d", lowest)
			return
		case low := <-lowChan:
			if lowest > low {
				lowest = low
			}
		}
	}

}

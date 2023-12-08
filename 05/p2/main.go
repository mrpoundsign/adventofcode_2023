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

	lowChan := make(chan int)
	var wg sync.WaitGroup

	// I know this is slow. Might work on it later.
	for i := 0; i < len(seedList.Seeds); i += 2 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			lowest := math.MaxInt

			ranges := [][2]int{
				{seedList.Seeds[i], seedList.Seeds[i+1] + seedList.Seeds[i] - 1},
			}

			for _, m := range seedMaps {
				ranges = m.LocationForRange(ranges)
			}

			for _, r := range ranges {
				if lowest > r[0] {
					lowest = r[0]
				}
			}
			lowChan <- lowest
		}(i)
	}

	wait := make(chan struct{})
	go func() {
		wg.Wait()
		close(wait)
	}()

	lowest := math.MaxInt

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

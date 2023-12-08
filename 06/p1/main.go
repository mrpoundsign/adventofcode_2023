package main

import (
	"bufio"
	"flag"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
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

	scanner := bufio.NewScanner(f)
	times := []int{}
	distances := []int{}

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		values := strings.Fields(line)
		if len(values) == 0 {
			continue
		}

		switch values[0] {
		case "Time:":
			for i := 1; i < len(values); i++ {
				time, err := strconv.Atoi(values[i])
				if err != nil {
					log.Fatal(err)
				}
				times = append(times, time)
			}
		case "Distance:":
			for i := 1; i < len(values); i++ {
				distance, err := strconv.Atoi(values[i])
				if err != nil {
					log.Fatal(err)
				}
				distances = append(distances, distance)
			}
		}
	}

	margin := 0

	for i, time := range times {
		wins := 0
		distance := distances[i]
		for j := 1; j <= time; j++ {
			speed := int(math.Pow(float64(j), 2)) / j

			dist := speed * (time - j)

			if dist > distance {
				// log.Println(time, j, dist, distance)
				wins++
			}
		}

		if margin == 0 {
			margin = wins
			continue
		}

		margin *= wins
	}

	log.Println("Margin of Error:", margin)
}

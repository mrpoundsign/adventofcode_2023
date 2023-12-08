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
	time := 0
	distance := 0

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
			nt, err := strconv.Atoi(strings.Join(values[1:], ""))
			if err != nil {
				log.Fatal(err)
			}
			time = nt
		case "Distance:":
			nd, err := strconv.Atoi(strings.Join(values[1:], ""))
			if err != nil {
				log.Fatal(err)
			}
			distance = nd
		}
	}

	wins := 0
	for j := 1; j <= time; j++ {
		speed := int(math.Pow(float64(j), 2)) / j

		dist := speed * (time - j)

		if dist > distance {
			// log.Println(time, j, dist, distance)
			wins++
		}
	}

	log.Println("Margin of Error:", wins)
}

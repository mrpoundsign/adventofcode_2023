package main

import (
	"bufio"
	"flag"
	"log"
	"os"
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

	res := 0
	for scanner.Scan() {
		s := -1
		e := -1
		line := scanner.Bytes()
		log.Printf("line: %q", line)
		for _, c := range line {
			if c < '0' || c > '9' {
				continue
			}

			if s == -1 {
				s = int(c - '0')
				continue
			}

			e = int(c - '0')
			log.Printf("s: %d e: %d", s, e)
		}

		if s == -1 {
			log.Fatalf("bad line: %q", line)
		}

		if e == -1 {
			e = s
		}
		log.Println(s*10+e, res)
		res += s*10 + e
	}
	log.Println(res)
}

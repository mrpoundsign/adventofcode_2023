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

	replacements := map[byte][]byte{
		'1': []byte("one"),
		'2': []byte("two"),
		'3': []byte("three"),
		'4': []byte("four"),
		'5': []byte("five"),
		'6': []byte("six"),
		'7': []byte("seven"),
		'8': []byte("eight"),
		'9': []byte("nine"),
	}

	res := 0
	for scanner.Scan() {
		lineS := scanner.Text()
		// if lineS == "abcone2threexyz" {
		// 	log.Fatal("fin")
		// }
		line := []byte(lineS)
		// log.Printf("line: %s", lineS)
		s := -1
		e := -1
		di := 0

	main:
		for i := 0; i < len(line); i++ {
			c := line[i]
			if c < '0' || c > '9' {
				// log.Printf("c: %q", c)
				found := false
				for ri, r := range replacements {
					rLen := len(r)
					if di >= rLen {
						continue
					}

					// if c == r[di] {
					if byteSliceEqual(r[:di+1], line[i-di:i+1]) {
						// log.Printf("c, r, di, rl: %q, %q, %q, %d", c, r[:di+1], line[cp-di:cp+1], di)
						if di == rLen-1 {
							i = i - di
							di = 0
							c = ri
							found = true
							// log.Println("found")
							break
						}

						di++
						continue main
					}
				}

				if !found {
					continue
				}
			}

			if s == -1 {
				s = int(c - '0')
				continue
			}

			e = int(c - '0')
		}

		if s == -1 {
			log.Fatalf("bad line: %q", line)
		}

		if e == -1 {
			e = s
		}
		res += s*10 + e
		log.Printf("line: %q, s: %d e: %d", line, s*10+e, res)
		// log.Println()
	}
	log.Println(res)
}

func byteSliceEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

package almanac

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strconv"
)

type ParseState int

const (
	STATE_SEED ParseState = iota
	STATE_SEED_LIST
	STATE_MAP
	STATE_MAP_LIST_SOURCE
	STATE_MAP_LIST_DESTINATION
	STATE_MAP_LIST_COUNT
)

func (ps ParseState) String() string {
	switch ps {
	case STATE_SEED:
		return "STATE_SEED"
	case STATE_SEED_LIST:
		return "STATE_SEED_LIST"
	case STATE_MAP:
		return "STATE_MAP"
	case STATE_MAP_LIST_SOURCE:
		return "STATE_MAP_LIST_SOURCE"
	case STATE_MAP_LIST_DESTINATION:
		return "STATE_MAP_LIST_DESTINATION"
	case STATE_MAP_LIST_COUNT:
		return "STATE_MAP_LIST_COUNT"
	default:
		return "STATE_UNKNOWN"
	}
}

type SeedList struct {
	Seeds []int
}

func (sl SeedList) String() string {
	var buf bytes.Buffer
	buf.WriteString("seeds:")
	for _, seed := range sl.Seeds {
		buf.WriteByte(' ')
		buf.WriteString(strconv.Itoa(seed))
	}
	return buf.String()
}

type SeedMap struct {
	From         string
	To           string
	Ranges       [][2]int
	Destinations []int
}

func (sm SeedMap) LocationFor(seed int) int {
	for i, r := range sm.Ranges {
		if seed < r[0] || seed > r[1] {
			continue
		}

		pos := seed - r[0]

		return sm.Destinations[i] + pos
	}
	return seed
}

func (sm SeedMap) locationForRange(i int, r [2]int) (error, [][2]int) {
	if len(sm.Ranges) <= i {
		return fmt.Errorf("no such range %d", i), nil
	}

	rm := sm.Ranges[i]
	newRanges := [][2]int{}
	rFirst := r[0]
	rLast := r[1]
	first := rm[0]
	last := rm[1]
	size := last - first

	if size < 0 {
		log.Panicln("zero size", first, last, rFirst, rLast)
	}

	if rLast < first || last < rFirst {
		// Not in this map
		return nil, [][2]int{}
	}

	log.Println("first", first, "last", last, "rfirst", rFirst, "rlast", rLast, "size", size)

	d := sm.Destinations[i]
	e := d + size

	if rFirst > first {
		log.Println("rFirst > first", rFirst, first)
		d += rFirst - first
		log.Println("d", d)
	}

	if rLast < last {
		log.Println("rLast < last", rLast, last)
		e -= last - rLast
		log.Println("e", e)
	}

	if rFirst < first {
		log.Println("rFirst < first", rFirst, first)
		count := first - rFirst
		newRanges = append(newRanges, [2]int{rFirst, rFirst + count})
	}

	if rLast > last {
		log.Println("rLast > last", rLast, last)
		count := rLast - last
		newRanges = append(newRanges, [2]int{rLast - count, rLast})
	}

	newRanges = append(newRanges, [2]int{d, e})
	return nil, newRanges
}

func (sm SeedMap) LocationForRange(ranges [][2]int) [][2]int {
	newRanges := [][2]int{}
	unprocessed := ranges
	for i, rm := range sm.Ranges {
		newUnprocessed := [][2]int{}
		for _, r := range unprocessed {
			rFirst := r[0]
			rLast := r[1]
			first := rm[0]
			last := rm[1]
			size := last - first

			if size < 0 {
				log.Panicln("zero size", first, last, rFirst, rLast)
			}

			if rLast < first || last < rFirst {
				newUnprocessed = append(newUnprocessed, r)
				// Not in this map
				continue
			}

			d := sm.Destinations[i]
			e := d + size

			if rFirst > first {
				d += rFirst - first
			}

			if rLast < last {
				e -= last - rLast
			}

			if rFirst < first {
				count := first - rFirst
				newUnprocessed = append(newUnprocessed, [2]int{rFirst, rFirst + count - 1})
			}

			if rLast > last {
				count := rLast - last
				newUnprocessed = append(newUnprocessed, [2]int{rLast - count + 1, rLast})
			}

			newRanges = append(newRanges, [2]int{d, e})
		}
		unprocessed = newUnprocessed
	}

	return append(unprocessed, newRanges...)
}

func (sm SeedMap) String() string {
	var buf bytes.Buffer
	buf.WriteString(sm.From)
	buf.WriteByte('-')
	buf.WriteString(sm.To)
	buf.WriteString(":\n")
	for _, source := range sm.Ranges {
		buf.WriteString(" s:")
		buf.WriteString(strconv.Itoa(source[0]))
		buf.WriteString(" ,:")
		buf.WriteString(strconv.Itoa(source[1]))
	}

	buf.WriteByte('\n')
	for _, dest := range sm.Destinations {
		buf.WriteString(" d:")
		buf.WriteString(strconv.Itoa(dest))
	}
	return buf.String()
}

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// Scan returns the next token and literal value.
func (p *Parser) Scan() (tok Token, lit string) {
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	p.buf.tok, p.buf.lit = p.s.Scan()
	return p.buf.tok, p.buf.lit
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }

func (p *Parser) Parse() (SeedList, []SeedMap) {
	var seeds SeedList
	var maps []SeedMap
	var currMap SeedMap
	source := -1
	dest := -1
	count := -1
	s := STATE_SEED
	inMap := false
	mapFrom := ""

	for {
		tok, lit := p.scan()
		switch tok {
		case NL:
			continue
		case WS:
			continue
		case EOF:
			if inMap {
				maps = append(maps, currMap)
			}
			return seeds, maps
		case ILLEGAL:
			log.Fatalln("ILLEGAL", lit)
		}

		switch s {
		case STATE_SEED:
			switch tok {
			case IDENT:
				if lit != "seeds" {
					log.Fatalf("expected 'seeds' got %s", lit)
				}
				continue
			case COLON:
				s = STATE_SEED_LIST
				seeds = SeedList{}
				continue
			}
		case STATE_SEED_LIST:
			switch tok {
			case NUMBER:
				seed, err := strconv.Atoi(lit)
				if err != nil {
					log.Fatalln("expected number")
				}
				seeds.Seeds = append(seeds.Seeds, seed)
				continue
			case IDENT:
				p.unscan()
				s = STATE_MAP
				continue
			}
		case STATE_MAP:
			// Continue scanning numbers
			switch tok {
			case NUMBER:
				if !inMap {
					break
				}

				s = STATE_MAP_LIST_DESTINATION
				p.unscan()
				continue
			case IDENT:
				// skip '-map' in map to/from
				if lit == "map" || lit == "to" {
					continue
				}
				if inMap {
					maps = append(maps, currMap)
					inMap = false
				}

				if mapFrom == "" {
					mapFrom = lit
					continue
				}

				currMap = SeedMap{
					From: mapFrom,
					To:   lit,
				}
				inMap = true
				continue
			case COLON:
				mapFrom = ""
				s = STATE_MAP_LIST_DESTINATION
				continue
			}
		case STATE_MAP_LIST_DESTINATION:
			switch tok {
			case IDENT:
				p.unscan()
				s = STATE_MAP
				continue
			case NL:
				if !inMap {
					log.Fatalln("STATE_MAP_LIST_DESTINATION when !inMap")
				}
				maps = append(maps, currMap)
				inMap = false
				s = STATE_MAP
				continue
			case NUMBER:
				d, err := strconv.Atoi(lit)
				if err != nil {
					log.Fatalln("expected number")
				}
				dest = d
				s = STATE_MAP_LIST_SOURCE
				continue
			}
		case STATE_MAP_LIST_SOURCE:
			switch tok {
			case NUMBER:
				sn, err := strconv.Atoi(lit)
				if err != nil {
					log.Fatalln("expected number")
				}
				source = sn
				s = STATE_MAP_LIST_COUNT
				continue
			}
		case STATE_MAP_LIST_COUNT:
			switch tok {
			case NUMBER:
				if count != -1 {
					log.Fatalf("expected only one number: %d", count)
				}
				count, err := strconv.Atoi(lit)
				if err != nil {
					log.Fatalln("expected number")
				}

				currMap.Ranges = append(currMap.Ranges, [2]int{source, source + count - 1})
				currMap.Destinations = append(currMap.Destinations, dest)

				s = STATE_MAP
				continue
			}
		}

		log.Fatalf("INVALID TOK %s(%d):%s STATE %s LINE %d", tok, tok, lit, s, p.s.Line+1)
	}
}

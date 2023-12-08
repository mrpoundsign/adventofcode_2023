package game

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"strconv"
)

const (
	VALUE_2 Card = iota
	VALUE_3
	VALUE_4
	VALUE_5
	VALUE_6
	VALUE_7
	VALUE_8
	VALUE_9
	VALUE_10
	VALUE_JACK
	VALUE_QUEEN
	VALUE_KING
	VALUE_ACE
)

const (
	STR_HIGH Strength = iota
	STR_PAIR
	STR_TWO_PAIR
	STR_THREE
	STR_FULL_HOUSE
	STR_FOUR
	STR_FIVE
)

const (
	CARD_2 = '2'
	CARD_3 = '3'
	CARD_4 = '4'
	CARD_5 = '5'
	CARD_6 = '6'
	CARD_7 = '7'
	CARD_8 = '8'
	CARD_9 = '9'
	CARD_T = 'T'
	CARD_J = 'J'
	CARD_Q = 'Q'
	CARD_K = 'K'
	CARD_A = 'A'
)

type Card int
type Strength int

func NewCard(c byte) (*Card, error) {
	card := Card(0)
	switch c {
	case CARD_2:
		card = VALUE_2
	case CARD_3:
		card = VALUE_3
	case CARD_4:
		card = VALUE_4
	case CARD_5:
		card = VALUE_5
	case CARD_6:
		card = VALUE_6
	case CARD_7:
		card = VALUE_7
	case CARD_8:
		card = VALUE_8
	case CARD_9:
		card = VALUE_9
	case CARD_T:
		card = VALUE_10
	case CARD_J:
		card = VALUE_JACK
	case CARD_Q:
		card = VALUE_QUEEN
	case CARD_K:
		card = VALUE_KING
	case CARD_A:
		card = VALUE_ACE
	default:
		return nil, fmt.Errorf("invalid card %c", c)
	}

	return &card, nil
}

func (c Card) String() string {
	switch c {
	case VALUE_2:
		return "2"
	case VALUE_3:
		return "3"
	case VALUE_4:
		return "4"
	case VALUE_5:
		return "5"
	case VALUE_6:
		return "6"
	case VALUE_7:
		return "7"
	case VALUE_8:
		return "8"
	case VALUE_9:
		return "9"
	case VALUE_10:
		return "10"
	case VALUE_JACK:
		return "J"
	case VALUE_QUEEN:
		return "Q"
	case VALUE_KING:
		return "K"
	case VALUE_ACE:
		return "A"
	}
	return "?"
}

type Hands []Hand

func (h Hands) Len() int           { return len(h) }
func (h Hands) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h Hands) Less(i, j int) bool { return h[i].Loser(h[j]) }

type Hand struct {
	cards       [5]Card
	cardMatches map[int]int
	bid         int
	str         Strength
}

func NewHand(cardsStr, bidStr string) *Hand {
	cards := []byte(cardsStr)

	if len(cards) != 5 {
		log.Fatalf("invalid hand %s", cardsStr)
	}

	bid, err := strconv.Atoi(bidStr)
	if err != nil {
		log.Fatalf("invalid bid %s", bidStr)
	}

	h := Hand{str: -1, bid: bid, cards: [5]Card{}}

	values := []int{}
	for i, c := range cards {
		card, err := NewCard(c)
		if err != nil {
			log.Fatalf("invalid card %c", c)
		}

		h.cards[i] = *card
		values = append(values, int(*card))
	}

	sort.Ints(values)

	h.cardMatches = make(map[int]int)
	curr := -1
	currNum := 0
	// log.Println("---", values)
	for _, v := range values {
		// log.Println(currNum, curr, v)
		if curr == -1 {
			curr = v
			currNum++
			continue
		}

		if v == curr {
			currNum++
			continue
		}

		if currNum > 1 {
			c, ok := h.cardMatches[currNum]
			if ok {
				h.cardMatches[currNum] = c + 1
			} else {
				h.cardMatches[currNum] = 1
			}
		}

		curr = v
		currNum = 1
	}

	// Add the last one
	if currNum > 1 {
		c, ok := h.cardMatches[currNum]
		if ok {
			h.cardMatches[currNum] = c + 1
		} else {
			h.cardMatches[currNum] = 1
		}
	}

	h.calculateStrength()

	return &h
}

func (h Hand) Strength() Strength {
	return h.str
}

func (h *Hand) calculateStrength() {

	h.str = STR_HIGH
	_, ok := h.cardMatches[5]
	if ok {
		h.str = STR_FIVE
		return
	}

	_, ok = h.cardMatches[4]
	if ok {
		h.str = STR_FOUR
		return
	}

	_, ok = h.cardMatches[3]
	if ok {
		_, ok = h.cardMatches[2]
		if ok {
			h.str = STR_FULL_HOUSE
			return
		}
		h.str = STR_THREE
		return
	}

	c, ok := h.cardMatches[2]
	if ok {
		if c == 2 {
			h.str = STR_TWO_PAIR
			return
		}
		h.str = STR_PAIR
		return
	}
}

func (h Hand) Loser(comp Hand) bool {
	if h.Strength() != comp.Strength() {
		return h.Strength() < comp.Strength()
	}

	for i := 0; i < 5; i++ {
		if h.cards[i] == comp.cards[i] {
			continue
		}

		return h.cards[i] < comp.cards[i]
	}

	return false
}

func (h Hand) Bid() int {
	return h.bid
}

func (h Hand) String() string {
	buff := bytes.Buffer{}
	for i := 0; i < 5; i++ {
		buff.WriteString(h.cards[i].String())
		buff.WriteString(" ")
	}

	for k, v := range h.cardMatches {
		buff.WriteString(fmt.Sprintf("%d:%c", v, k))
	}
	return fmt.Sprintf("(%s):%d:%d", buff.String(), h.Strength(), h.bid)
}

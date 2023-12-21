package main

import (
	"fmt"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

var fn = "./day_07/input.txt"

type hand struct {
	value string
	bid   int
	count int
	score []int
}

func newHandA(s string, b int) *hand {
	h := &hand{
		value: s,
		bid:   b,
	}

	l := utf8.RuneCount([]byte(s))
	rc := make(map[rune]int, l)
	h.score = make([]int, l)

	for i, r := range s {
		rc[r]++

		switch r {
		case 'A':
			h.score[i] = 14
		case '2':
			h.score[i] = 2
		case '3':
			h.score[i] = 3
		case '4':
			h.score[i] = 4
		case '5':
			h.score[i] = 5
		case '6':
			h.score[i] = 6
		case '7':
			h.score[i] = 7
		case '8':
			h.score[i] = 8
		case '9':
			h.score[i] = 9
		case 'T':
			h.score[i] = 10
		case 'J':
			h.score[i] = 11
		case 'Q':
			h.score[i] = 12
		case 'K':
			h.score[i] = 13
		default:
			panic("failed to parse character")
		}

	}

	rr := make([]int, 0, len(rc))

	for _, c := range rc {
		rr = append(rr, c)
	}

	slices.Sort(rr)
	slices.Reverse(rr)

	switch {
	case rr[0] == 5:
		h.count = 6
	case rr[0] == 4:
		h.count = 5
	case rr[0] == 3 && rr[1] == 2:
		h.count = 4
	case rr[0] == 3:
		h.count = 3
	case rr[0] == 2 && rr[1] == 2:
		h.count = 2
	case rr[0] == 2:
		h.count = 1
	}

	return h
}

func newHandB(s string, b int) *hand {
	h := &hand{
		value: s,
		bid:   b,
	}

	l := utf8.RuneCount([]byte(s))
	rc := make(map[rune]int, l)
	h.score = make([]int, l)

	var jc int

	for i, r := range s {
		if r != 'J' {
			rc[r]++
		} else {
			if _, ok := rc[r]; !ok {
				rc[r] = 0
			}
		}

		switch r {
		case 'A':
			h.score[i] = 14
		case '2':
			h.score[i] = 2
		case '3':
			h.score[i] = 3
		case '4':
			h.score[i] = 4
		case '5':
			h.score[i] = 5
		case '6':
			h.score[i] = 6
		case '7':
			h.score[i] = 7
		case '8':
			h.score[i] = 8
		case '9':
			h.score[i] = 9
		case 'T':
			h.score[i] = 10
		case 'J':
			jc++
			h.score[i] = 1
		case 'Q':
			h.score[i] = 12
		case 'K':
			h.score[i] = 13
		default:
			panic("failed to parse character")
		}

	}

	rr := make([]int, 0, len(rc))

	for _, c := range rc {
		rr = append(rr, c)
	}

	slices.Sort(rr)
	slices.Reverse(rr)

	switch {
	case rr[0]+jc == 5:
		h.count = 6
	case rr[0]+jc == 4:
		h.count = 5
	case rr[0]+jc == 3 && rr[1] == 2:
		h.count = 4
	case rr[0]+jc == 3:
		h.count = 3
	case rr[0] == 2 && (rr[1] == 2 || jc == 2):
		h.count = 2
	case rr[0]+jc == 2:
		h.count = 1
	}

	return h
}

type hands []*hand

func (hh hands) Len() int {
	return len(hh)
}

func (hh hands) Swap(i, j int) {
	hh[i], hh[j] = hh[j], hh[i]
}

func (hh hands) Less(i, j int) bool {
	a, b := hh[i], hh[j]

	if a.count != b.count {
		return a.count < b.count
	}

	for k := 0; k < len(a.score); k++ {
		c, d := a.score[k], b.score[k]
		if c != d {
			return c < d
		}
	}

	return false
}

func parse(s string) ([]string, []int) {
	bb := strings.Split(s, "\n")
	vv := make([]string, 0, len(bb))
	cc := make([]int, 0, len(bb))

	for _, b := range bb {
		if b == "" {
			continue
		}

		h := strings.Split(b, " ")
		c, err := strconv.Atoi(h[1])
		if err != nil {
			panic(err)
		}

		vv = append(vv, h[0])
		cc = append(cc, c)
	}

	return vv, cc
}

func makeCalc(p func(string, int) *hand) func([]string, []int) int {
	return func(vv []string, cc []int) int {
		aa := make(hands, len(vv))

		for i, v := range vv {
			aa[i] = p(v, cc[i])
		}

		sort.Sort(aa)

		var a int
		for i, h := range aa {
			a += h.bid * (i + 1)
		}

		return a
	}
}

func main() {
	bb, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}

	vv, cc := parse(string(bb))
	a := makeCalc(newHandA)(vv, cc)
	b := makeCalc(newHandB)(vv, cc)

	fmt.Println("part one value: ", a)
	fmt.Println("part two value: ", b)
}

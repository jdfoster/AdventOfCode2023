package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	fn = "./day_03/input.txt"
	re = regexp.MustCompile(`\d+`)
)

func main() {
	bb, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}

	s := string(bb)
	ll := strings.IndexRune(s, '\n') + 1

	vv := make([]int, 0, 8)
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if i != 0 || j != 0 {
				vv = append(vv, ll*i+j)
			}
		}
	}

	var total int
	gs := make(map[int][]int)

	for _, p := range re.FindAllIndex(bb, -1) {
		ai, bi := p[0], p[1]
		ss := string(bb[ai:bi])
		n, err := strconv.Atoi(ss)
		if err != nil {
			panic(err)
		}

	LOOP:
		for i := ai; i < bi; i++ {
			for _, v := range vv {
				if j := i + v; j >= 0 && j < len(bb) {
					r := rune(bb[j])

					if !unicode.IsDigit(r) && r != '.' && r != '\n' {
						total += n

						if r == '*' {
							gg := gs[j]
							gs[j] = append(gg, n)
						}

						break LOOP
					}
				}
			}
		}
	}

	fmt.Println("part one value: ", total)

	var product int
	for _, gg := range gs {
		if len(gg) == 2 {
			product += (gg[0] * gg[1])
		}
	}

	fmt.Println("part two value: ", product)
}

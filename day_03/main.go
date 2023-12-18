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

						break LOOP
					}
				}
			}
		}
	}

	fmt.Println("part one value: ", total)
}

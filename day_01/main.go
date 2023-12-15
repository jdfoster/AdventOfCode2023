package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"unicode"
)

var fn = "./day_01/input.txt"

func main() {
	f, err := os.Open(fn)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	s := bufio.NewScanner(f)

	var sum int

	for s.Scan() {

		var (
			a, b  rune
			first bool
		)

		rr := s.Text()

		for _, r := range rr {
			if unicode.IsNumber(r) {
				b = r

				if !first {
					a = r
					first = true
				}
			}
		}

		v, err := strconv.Atoi(string([]rune{a, b}))
		if err != nil {
			panic(err)
		}

		sum += v
	}

	fmt.Println("part one value: ", sum)
}

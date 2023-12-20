package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var fn = "./day_06/input.txt"

func parseA(s string) ([]int, []int) {
	r := make([][]int, 2)
	ss := strings.Split(s, "\n")

	for i := 0; i < 2; i++ {
		vv := strings.Fields(ss[i])[1:]
		r[i] = make([]int, len(vv))
		for j, v := range vv {
			var err error
			r[i][j], err = strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
		}
	}

	return r[0], r[1]
}

func parseB(s string) (int, int) {
	r := make([]int, 2)
	ss := strings.Split(s, "\n")

	for i := 0; i < 2; i++ {
		b := []byte(ss[i])
		bb := make([]rune, 0, utf8.RuneCount(b))
		for _, c := range ss[i] {
			if unicode.IsDigit(c) {
				bb = append(bb, c)
			}
		}

		var err error
		r[i], err = strconv.Atoi(string(bb))
		if err != nil {
			panic(err)
		}
	}

	return r[0], r[1]
}

func solve(t, d, i int) bool {
	dt := t - i
	x := dt * i
	return x > d
}

func find(t, d int) (int, int) {
	var a, b int

	for i := d / t; i < t; i++ {
		if solve(t, d, i) {
			a = i
			break
		}
	}

	for j := t; j > 0; j-- {
		if solve(t, d, j) {
			b = j
			break
		}
	}

	return a, b
}

func main() {
	bb, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}

	iat, iad := parseA(string(bb))

	a := 1

	for i := 0; i < len(iat); i++ {
		x, y := find(iat[i], iad[i])
		a *= (y - x) + 1
	}

	fmt.Println("part one value: ", a)

	ibt, ibd := parseB(string(bb))
	bx, by := find(ibt, ibd)
	b := (by - bx) + 1

	fmt.Println("part two value: ", b)
}

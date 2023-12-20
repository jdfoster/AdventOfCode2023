package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var fn = "./day_06/input.txt"

func parse(s string) ([]int, []int) {
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

func solve(t, d int) []int {
	rr := make([]int, 0, t)

	for i := 0; i < t; i++ {
		dt := t - i
		x := dt * i

		if x >= d {
			rr = append(rr, i)
		}
	}

	return rr
}

func main() {
	bb, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}

	ia, ib := parse(string(bb))

	a := 1

	for i := 0; i < len(ia); i++ {
		a *= len(solve(ia[i], ib[i]))
	}

	fmt.Println("part one value: ", a)
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var fn = "./day_04/input.txt"

func parseGame(s string) ([]int, []int) {
	ss := strings.Split(s, ":")
	ll := strings.Split(ss[1], "|")

	r := make([][]int, len(ll))

	for i, l := range ll {
		gg := strings.Split(strings.TrimSpace(l), " ")
		r[i] = make([]int, 0, len(gg))

		for _, g := range gg {
			if len(g) > 0 {
				n, err := strconv.Atoi(g)
				if err != nil {
					panic(err)
				}

				r[i] = append(r[i], n)
			}
		}
	}

	return r[0], r[1]
}

func findMatched(aa, bb []int) []int {
	as := make(map[int]struct{}, len(aa))
	for _, a := range aa {
		as[a] = struct{}{}
	}

	rr := make([]int, 0, len(bb))
	for _, b := range bb {
		if _, ok := as[b]; ok {
			rr = append(rr, b)
		}
	}

	return rr
}

func calcScore(a []int) int {
	var v int
	l := len(a)

	if l > 0 {
		v = 1

	}

	for i := 0; i < l -1; i++ {
		v *= 2
	}


	fmt.Println(l, v)

	return v
}

func main() {
	f, err := os.Open(fn)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	s := bufio.NewScanner(f)

	var a int

	for s.Scan() {
		a += calcScore(findMatched(parseGame(s.Text())))
	}

	fmt.Println("part one value: ", a)
}

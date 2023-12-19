package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var fn = "./day_05/input.txt"

type rule = func(int) (int, bool)

func mkRule(a, b, c int) rule {
	return func(s int) (int, bool) {
		d := s - b

		if d > -1 && d <= c {
			return a + d, true
		}

		return 0, false
	}
}

type lookup = func(int) int

func mkRules(rr ...rule) lookup {
	return func(s int) int {
		for _, r := range rr {
			if v, ok := r(s); ok {
				return v
			}
		}
		return s
	}
}

func mkPath(ss ...lookup) lookup {
	return func(i int) int {
		for _, s := range ss {
			i = s(i)
		}

		return i
	}
}

func parse(b []byte) ([]int, [][][]int) {
	bb := strings.Split(string(b), "\n\n")
	sa := strings.Split(strings.TrimPrefix(bb[0], "seeds: "), " ")
	sb := bb[1:]

	seeds := make([]int, len(sa))
	rules := make([][][]int, len(sb))
	var err error

	for i, v := range sa {
		seeds[i], err = strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
	}

	for i, si := range sb {
		rr := strings.Split(si, "\n")[1:]

		if l := len(rr) - 1; rr[l] == "" {
			rr = rr[:l]
		}

		rules[i] = make([][]int, len(rr))
		for j, r := range rr {
			vv := strings.Split(r, " ")
			rules[i][j] = make([]int, len(vv))
			for k, v := range vv {
				if v == "" {
					continue
				}

				rules[i][j][k], err = strconv.Atoi(v)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	return seeds, rules
}

func main() {
	bb, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}

	ss, rrr := parse(bb)

	ggg := make([]lookup, len(rrr))
	for i, rr := range rrr {
		gg := make([]rule, len(rr))
		for j, r := range rr {
			if len(r) != 3 {
				panic("rule is not 3 elements long")
			}
			gg[j] = mkRule(r[0], r[1], r[2])
		}
		ggg[i] = mkRules(gg...)
	}

	p := mkPath(ggg...)

	vv := make([]int, len(ss))
	for i, s := range ss {
		vv[i] = p(s)
	}

	m := math.MaxInt
	for _, v := range vv {
		if v < m {
			m = v
		}
	}

	fmt.Println("part one value: ", m)
}

package main

import (
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
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

func unzip(e int, xx ...int) []int {
	vv := make([]int, 0, len(xx))

	for i, x := range xx {
		if i%2 == e {
			vv = append(vv, x)
		}
	}

	return vv
}

func main() {
	bb, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}

	ia, rrr := parse(bb)

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

	pt := mkPath(ggg...)
	proc := func(xx ...int) []int {
		v := make([]int, len(xx))
		for i, x := range xx {
			v[i] = pt(x)
		}

		return v
	}

	fmt.Println("part one value: ", slices.Min(proc(ia...)))

	iba := unzip(0, ia...)
	ibb := unzip(1, ia...)

	if len(iba) != len(ibb) {
		panic("lengths are not equal")
	}

	var wg sync.WaitGroup
	wg.Add(len(iba))
	rb := make([]int, len(iba))

	for i := 0; i < len(iba); i++ {
		go func(j int) {
			rb[j] = math.MaxInt
			a, b := iba[j], ibb[j]

			for k := 0; k < b; k++ {
				rb[j] = min(rb[j], pt(a+k))
			}

			fmt.Println(fmt.Sprintf(">>> [%d] %d", j, rb[j]))
			wg.Done()
		}(i)
	}

	wg.Wait()

	fmt.Println("part two value: ", slices.Min(rb))
}

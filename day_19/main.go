package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	fn         = "./day_19/input.txt"
	reInstruct = regexp.MustCompile(`([a-z]{2,3})\{(([a-z][><][0-9]*:[a-zA-Z]*,?)+)([a-zA-Z]*)\}`)
)

type outcome uint8

const (
	REJECT outcome = iota
	ACCEPT
)

func runeByIndex(s string, i int) rune {
	for j, r := range s {
		if j == i {
			return r
		}
	}

	panic(fmt.Sprintf("failed to find rune with index %d in string %q", i, s))
}

func stringByIndex(s string, i, j int) string {
	var b strings.Builder

	for k, r := range s {
		if k > j {
			break
		}
		if k >= i {
			b.WriteRune(r)
		}
	}

	return b.String()
}

type shape map[rune]int
type rule func(s ruleLookup, m shape) outcome
type ruleLookup map[string]rule

func makeDestination(d string) rule {
	return func(s ruleLookup, m shape) outcome {
		switch d {
		case "A":
			return ACCEPT
		case "R":
			return REJECT
		}

		r, ok := s[d]
		if !ok {
			panic(fmt.Sprintf("failed to find label in state: %s", d))
		}

		return r(s, m)
	}
}

func makeRule(c, o rune, v int, d string) func(rule) rule {
	return func(next rule) rule {
		return func(s ruleLookup, m shape) outcome {
			x, ok := m[c]
			if !ok {
				panic(fmt.Sprintf("failed to find key in element: %s", string(c)))
			}

			var p bool
			switch o {
			case '>':
				p = x > v
			case '<':
				p = x < v
			default:
				panic(fmt.Sprintf("unknown operator value: %s", string(o)))
			}

			if p {
				return makeDestination(d)(s, m)
			}

			return next(s, m)
		}
	}
}

func makeFilter(gg []string, d string) rule {
	rr := make([]func(rule) rule, len(gg))

	for i, g := range gg {
		c := runeByIndex(g, 0)
		o := runeByIndex(g, 1)
		j := strings.IndexRune(g, ':')
		k := utf8.RuneCountInString(g)
		d := stringByIndex(g, j+1, k)
		vi := stringByIndex(g, 2, j-1)
		v, err := strconv.Atoi(vi)
		if err != nil {
			panic(err)
		}

		rr[i] = makeRule(c, o, v, d)
	}

	slices.Reverse(rr)

	var r rule
	for i, ri := range rr {
		if i == 0 {
			r = ri(makeDestination(d))
			continue
		}

		r = ri(r)
	}

	return r
}

func parse(s string) (ruleLookup, []shape) {
	part := strings.Split(s, "\n\n")
	qq := strings.Split(part[0], "\n")
	rr := make(ruleLookup, len(qq))

	for _, q := range qq {
		mm := reInstruct.FindAllStringSubmatch(q, -1)
		gg := strings.Split(strings.TrimSuffix(mm[0][2], ","), ",")
		h := mm[0][4]
		l := mm[0][1]
		rr[l] = makeFilter(gg, h)
	}

	kk := strings.Split(part[1], "\n")
	ee := make([]shape, len(kk))

	for i, k := range kk {
		ee[i] = make(map[rune]int, 3)
		hh := strings.Split(strings.TrimSuffix(strings.TrimPrefix(k, "{"), "}"), ",")

		for _, h := range hh {
			if h == "" {
				continue
			}

			r := runeByIndex(h, 0)
			k := utf8.RuneCountInString(h)
			vi := stringByIndex(h, 2, k)
			v, err := strconv.Atoi(vi)
			if err != nil {
				panic(err)
			}

			ee[i][r] = v
		}

	}

	return rr, ee[:len(ee)-1]
}

func main() {
	bb, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}

	rr, ss := parse(string(bb))

	var a int

	for _, s := range ss {
		r, ok := rr["in"]
		if !ok {
			panic("missing starting rule")
		}
		if r(rr, s) == ACCEPT {
			var sum int
			for _, v := range s {
				sum += v
			}
			a += sum
			// fmt.Println(i, sum, a)
		}
	}

	fmt.Println("part one value: ", a)
}

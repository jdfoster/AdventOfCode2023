package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	fn         = "./day_19/input.txt"
	reInstruct = regexp.MustCompile(`([a-z]{2,3})\{(([a-z][><][0-9]*:[a-zA-Z]*,?)+)([a-zA-Z]*)\}`)
)

func stringByIndex(s string, i, j int) string {
	var b strings.Builder

	for k, r := range s {
		if k >= j {
			break
		}
		if k >= i {
			b.WriteRune(r)
		}
	}

	return b.String()
}

type op uint8

const (
	LESS_THAN op = iota
	GREATER_THAN
)

type criterion struct {
	component string
	operator  op
	value     int
	whenTrue  string
}

type instruct struct {
	label     string
	criteria  []*criterion
	whenFalse string
}

func (ins instruct) findDestination(m shape) string {
	for _, c := range ins.criteria {
		vi, ok := m[c.component]
		if !ok {
			panic(fmt.Sprintf("failed to find component %q in shape", c.component))
		}

		var r bool
		switch c.operator {
		case LESS_THAN:
			r = vi < c.value
		case GREATER_THAN:
			r = vi > c.value
		}

		if r {
			return c.whenTrue
		}
	}

	return ins.whenFalse
}

func newInstruct(s string) *instruct {
	aa := reInstruct.FindAllStringSubmatch(s, -1)
	rr := strings.Split(strings.TrimSuffix(aa[0][2], ","), ",")

	r := &instruct{
		label:     aa[0][1],
		criteria:  make([]*criterion, len(rr)),
		whenFalse: aa[0][4],
	}

	for i, ri := range rr {
		j := strings.IndexRune(ri, ':')
		k := utf8.RuneCountInString(ri)
		vi := stringByIndex(ri, 2, j)

		r.criteria[i] = &criterion{
			component: stringByIndex(ri, 0, 1),
			whenTrue:  stringByIndex(ri, j+1, k),
		}

		switch op := stringByIndex(ri, 1, 2); op {
		case "<":
			r.criteria[i].operator = LESS_THAN
		case ">":
			r.criteria[i].operator = GREATER_THAN
		default:
			panic(fmt.Sprintf("unknown operator %q", op))
		}

		var err error
		r.criteria[i].value, err = strconv.Atoi(vi)
		if err != nil {
			panic(err)
		}
	}

	return r
}

type outcome uint8

const (
	REJECT outcome = iota
	ACCEPT
)

type lookup map[string]*instruct
type shape map[string]int

func (u lookup) validateShape(s shape, l string) outcome {
	ins, ok := u[l]
	if !ok {
		panic(fmt.Sprintf("failed to find label %q", l))
	}

	d := ins.findDestination(s)

	switch d {
	case "A":
		return ACCEPT
	case "R":
		return REJECT
	}

	return u.validateShape(s, d)
}

func newLookup(ss []string) lookup {
	ll := make(map[string]*instruct, len(ss))

	for _, s := range ss {
		r := newInstruct(s)
		ll[r.label] = r
	}

	return ll
}

func parse(f string) (lookup, []shape) {
	part := strings.Split(f, "\n\n")
	qq := strings.Split(part[0], "\n")
	rr := newLookup(qq)

	kk := strings.Split(part[1], "\n")
	ss := make([]shape, len(kk))

	var n int
	for i, k := range kk {
		ss[i] = make(map[string]int, 3)
		cc := strings.Split(strings.TrimSuffix(strings.TrimPrefix(k, "{"), "}"), ",")

		for _, c := range cc {
			if c == "" {
				n++
				continue
			}

			vi := strings.Split(c, "=")
			k, v := vi[0], vi[1]

			var err error
			ss[i][k], err = strconv.Atoi(v)
			if err != nil {
				panic(err)
			}
		}
	}

	return rr, ss[:len(ss)-n]
}

func main() {
	bb, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}

	rr, ss := parse(string(bb))

	var a int

	for _, s := range ss {
		if rr.validateShape(s, "in") == ACCEPT {
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

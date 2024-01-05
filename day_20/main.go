package main

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

var fn = "./day_20/input.txt"

type polarity uint8

const (
	LOW polarity = iota
	HIGH
)

type pulse struct {
	src   string
	dest  string
	kind  polarity
	count int
}

type part interface {
	process(p pulse) []pulse
}

type flipFlop struct {
	label   string
	enabled bool
	targets []string
}

func (f *flipFlop) process(p pulse) []pulse {
	if p.kind == LOW {
		f.enabled = !f.enabled

		var k polarity

		if f.enabled {
			k = HIGH
		}

		rr := make([]pulse, len(f.targets))

		for i, d := range f.targets {
			rr[i] = pulse{
				src:   f.label,
				dest:  d,
				kind:  k,
				count: p.count,
			}
		}

		return rr
	}

	return nil
}

type conjunction struct {
	label   string
	sources map[string]polarity
	targets []string
	cycles  map[string]int
}

func (c *conjunction) intervals() ([]int, bool) {
	if len(c.sources) != len(c.cycles) {
		return nil, false
	}

	r := make([]int, 0, len(c.cycles))
	for _, v := range c.cycles {
		r = append(r, v)
	}

	return r, true
}

func (c *conjunction) equal(p polarity) bool {
	for _, v := range c.sources {
		if v != p {
			return false
		}
	}

	return true
}

func (c *conjunction) process(p pulse) []pulse {
	_, ok := c.sources[p.src]
	if !ok {
		panic(fmt.Sprintf("component %q: failed to find source %q", c.label, p.src))
	}

	c.sources[p.src] = p.kind

	if p.kind == HIGH {
		if _, ok := c.cycles[p.src]; !ok {
			c.cycles[p.src] = p.count
			// fmt.Println(p.src, p.count)
		}
	}

	var k polarity
	if !c.equal(HIGH) {
		k = HIGH
	}

	rr := make([]pulse, len(c.targets))
	for i, d := range c.targets {
		rr[i] = pulse{
			src:   c.label,
			dest:  d,
			kind:  k,
			count: p.count,
		}
	}

	return rr
}

type boardcaster struct {
	label   string
	targets []string
}

func (b *boardcaster) process(p pulse) []pulse {
	rr := make([]pulse, len(b.targets))

	for i, d := range b.targets {
		rr[i] = pulse{
			src:   b.label,
			dest:  d,
			count: p.count,
		}
	}

	return rr
}

type machine struct {
	items    map[string]part
	count    int
	terminal *conjunction
}

func (m *machine) findIntervals() []int {
LOOP:
	v, ok := m.terminal.intervals()

	if !ok {
		m.button()
		goto LOOP
	}

	return v
}

func (m *machine) run(v []pulse) (int, int) {
	var loop func(qq [][]pulse) (int, int)

	loop = func(qq [][]pulse) (int, int) {
		var countLow, countHigh int

		if len(qq) > 0 {
			rr := make([][]pulse, 0, len(qq))

			for _, pp := range qq {
				if pp == nil {
					continue
				}

				for _, p := range pp {
					switch p.kind {
					case LOW:
						countLow++
					case HIGH:
						countHigh++
					}

					c, ok := m.items[p.dest]
					if !ok {
						continue
					}

					rr = append(rr, c.process(p))
				}
			}

			a, b := loop(rr)

			countLow += a
			countHigh += b
		}

		return countLow, countHigh
	}

	return loop([][]pulse{v})
}

func (m *machine) button() (int, int) {
	m.count++

	s, ok := m.items["broadcaster"]
	if !ok {
		panic("failed to find broadcaster module")
	}

	a, b := m.run(s.process(pulse{count: m.count}))

	return a + 1, b
}

func (m *machine) rep(n int) int {
	var a, b int

	for i := 0; i < n; i++ {
		ai, bi := m.button()
		a += ai
		b += bi
	}

	return a * b
}

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

type component uint8

const (
	FLIP_FLOP component = iota
	CONJUNCTION
	BOARDCASTER
)

func parseLabel(l string) (r string, k component) {
	switch {
	case strings.IndexRune(l, '%') == 0:
		r = stringByIndex(l, 1, utf8.RuneCountInString(l))
		k = FLIP_FLOP
	case strings.IndexRune(l, '&') == 0:
		r = stringByIndex(l, 1, utf8.RuneCountInString(l))
		k = CONJUNCTION
	case l == "broadcaster":
		r = l
		k = BOARDCASTER
	default:
		panic(fmt.Sprintf("failed to parse label %q", l))
	}

	return
}

func parse(p string) *machine {
	ss := strings.Split(p, "\n")
	m := &machine{
		items: make(map[string]part, len(ss)),
	}
	srcs := make(map[string]map[string]polarity, len(ss))
	cons := make(map[string]*conjunction, len(ss))

	for _, s := range ss {
		if s == "" {
			continue
		}

		pp := strings.Split(s, " -> ")
		dd := strings.Split(pp[1], ",")

		targets := make([]string, len(dd))
		label, kind := parseLabel(pp[0])

		for i, d := range dd {
			di := strings.TrimSpace(d)
			targets[i] = di

			if _, ok := srcs[di]; !ok {
				srcs[di] = make(map[string]polarity, 5)
			}

			srcs[di][label] = LOW
		}

		switch kind {
		case FLIP_FLOP:
			m.items[label] = &flipFlop{
				label:   label,
				targets: targets,
			}
		case CONJUNCTION:
			if _, ok := srcs[label]; !ok {
				srcs[label] = make(map[string]polarity, 5)
			}

			cons[label] = &conjunction{
				label:   label,
				sources: srcs[label],
				targets: targets,
				cycles:  make(map[string]int, 5),
			}

			m.items[label] = cons[label]
		case BOARDCASTER:
			m.items[label] = &boardcaster{
				label:   label,
				targets: targets,
			}
		}
	}

	// df is the module that connect to rx
	m.terminal = cons["df"]

	return m
}

// https://siongui.github.io/2017/06/03/go-find-lcm-by-gcd/
// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int, integers ...int) int {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

func main() {
	bb, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}

	m := parse(string(bb))

	fmt.Println("part one value: ", m.rep(1000))

	bi := m.findIntervals()
	fmt.Println("part two value: ", LCM(bi[0], bi[1], bi[2:]...))
}

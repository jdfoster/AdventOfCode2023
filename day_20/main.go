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
	src  string
	dest string
	kind polarity
}

type component interface {
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
				src:  f.label,
				dest: d,
				kind: k,
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
}

func (c *conjunction) process(p pulse) []pulse {
	_, ok := c.sources[p.src]
	if !ok {
		panic(fmt.Sprintf("component %q: failed to find source %q", c.label, p.src))
	}

	c.sources[p.src] = p.kind

	var allHigh bool = true
	for _, v := range c.sources {
		allHigh = allHigh && v == HIGH
		if !allHigh {
			break
		}
	}

	var k polarity
	if !allHigh {
		k = HIGH
	}
	rr := make([]pulse, len(c.targets))

	for i, d := range c.targets {
		rr[i] = pulse{
			src:  c.label,
			dest: d,
			kind: k,
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
			src:  b.label,
			dest: d,
		}
	}

	return rr
}

type machine struct {
	items map[string]component
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
	s, ok := m.items["broadcaster"]
	if !ok {
		panic("failed to find broadcaster module")
	}

	a, b := m.run(s.process(pulse{}))

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

func parseLabel(l string) (r string, k uint8) {
	switch {
	case strings.IndexRune(l, '%') == 0:
		r = stringByIndex(l, 1, utf8.RuneCountInString(l))
		k = 0
	case strings.IndexRune(l, '&') == 0:
		r = stringByIndex(l, 1, utf8.RuneCountInString(l))
		k = 1
	case l == "broadcaster":
		r = l
		k = 2
	default:
		k = 3
	}

	return
}

func parse(p string) *machine {
	ss := strings.Split(p, "\n")
	m := &machine{items: make(map[string]component, len(ss))}
	srcs := make(map[string]map[string]polarity, len(ss))

	for _, s := range ss {
		if s == "" {
			continue
		}

		pp := strings.Split(s, " -> ")
		dd := strings.Split(pp[1], ",")

		targets := make([]string, len(dd))
		label, k := parseLabel(pp[0])

		for i, d := range dd {
			di := strings.TrimSpace(d)
			targets[i] = di

			if _, ok := srcs[di]; !ok {
				srcs[di] = make(map[string]polarity, 5)
			}

			srcs[di][label] = LOW
		}

		switch k {
		case 0:
			m.items[label] = &flipFlop{
				label:   label,
				targets: targets,
			}
		case 1:
			if _, ok := srcs[label]; !ok {
				srcs[label] = make(map[string]polarity, 5)
			}

			m.items[label] = &conjunction{
				label:   label,
				sources: srcs[label],
				targets: targets,
			}
		case 2:
			m.items[label] = &boardcaster{
				label:   label,
				targets: targets,
			}
		default:
			panic(fmt.Sprintf("failed to parse label %q", pp[0]))
		}

	}

	return m
}

func main() {
	bb, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}

	m := parse(string(bb))
	a := m.rep(1000)

	fmt.Println("part one value: ", a)
}

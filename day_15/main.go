package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

var fn = "./day_15/input.txt"

func hash(s string) int {
	var h int

	for _, r := range s {
		h += int(r)
		h *= 17
		h = h % 256
	}

	return h
}

func isEqual(ss string) (int, bool) {
	if strings.ContainsRune(ss, '=') {
		var b strings.Builder
		for _, s := range ss {
			if unicode.IsDigit(s) {
				b.WriteRune(s)
			}
		}

		v, err := strconv.Atoi(b.String())
		if err != nil {
			panic(err)
		}

		return v, true
	}

	return 0, false
}

func extract(ss string) (string, int, bool) {
	var b strings.Builder
	for _, s := range ss {
		if unicode.IsLetter(s) {
			b.WriteRune(s)
			continue
		}

		break
	}

	f, ok := isEqual(ss)
	return b.String(), f, ok
}

type instruction struct {
	value   string
	label   string
	focal   int
	isEqual bool
}

func newInstruction(ss string) *instruction {
	e := &instruction{
		value: ss,
	}
	e.label, e.focal, e.isEqual = extract(ss)

	return e
}

func parse(s string) []*instruction {
	si := strings.Split(s, "\n")
	vv := strings.Split(si[0], ",")
	r := make([]*instruction, len(vv))

	for i, v := range vv {
		r[i] = newInstruction(v)
	}

	return r
}

type box struct {
	focal map[string]int
	order []string
}

func (b *box) add(l string, f int) {
	b.focal[l] = f

	if !slices.Contains(b.order, l) {
		b.order = append(b.order, l)
	}
}

func (b *box) remove(l string) {
	delete(b.focal, l)

	if i := slices.Index(b.order, l); i > -1 {
		b.order = slices.Delete(b.order, i, i+1)
	}
}

func (b box) calc() int {
	var r int
	for i, l := range b.order {
		f, ok := b.focal[l]
		if !ok {
			panic(fmt.Sprintf("label %q not found", l))
		}

		r += (i + 1) * f
	}

	return r
}

func newBox() *box {
	return &box{
		focal: map[string]int{},
	}
}

type boxes []*box

func (bb boxes) calc() int {
	var r int
	for i, b := range bb {
		r += (i + 1) * b.calc()
	}

	return r
}

func newBoxes() boxes {
	rr := make(boxes, 256)

	for i := 0; i < len(rr); i++ {
		rr[i] = newBox()
	}

	return rr
}

func main() {
	b, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}

	ss := parse(string(b))
	bb := newBoxes()

	var a int

	for _, s := range ss {
		a += hash(s.value)

		l := hash(s.label)
		if s.isEqual {
			bb[l].add(s.label, s.focal)
		} else {
			bb[l].remove(s.label)
		}
	}

	fmt.Println("part one value: ", a)
	fmt.Println("part two value: ", bb.calc())
}

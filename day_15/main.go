package main

import (
	"fmt"
	"os"
	"strings"
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

func parse(s string) []string {
	si := strings.Split(s, "\n")[0]
	return strings.Split(si, ",")
}

func main() {
	bb, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}

	ss := parse(string(bb))

	var a int

	for _, s := range ss {
		a += hash(s)
	}

	fmt.Println("part one value: ", a)
}

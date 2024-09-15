package main

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"strings"
)

type printFmt uint8

const (
	printFmtNone printFmt = iota
	printFmtShort
	printFmtLong
	printFmtSpace
	printFmtHeads
)

func (f *printFmt) UnmarshalText(b []byte) error {
	s := string(b)
	switch s {
	case "none":
		*f = printFmtNone
	case "short":
		*f = printFmtShort
	case "long":
		*f = printFmtLong
	case "space":
		*f = printFmtSpace
	case "heads":
		*f = printFmtHeads
	default:
		return fmt.Errorf("invalid format: %s", b)
	}
	return nil
}

func (f printFmt) MarshalText() ([]byte, error) {
	switch f {
	case printFmtNone:
		return []byte("none"), nil
	case printFmtShort:
		return []byte("short"), nil
	case printFmtLong:
		return []byte("long"), nil
	case printFmtSpace:
		return []byte("space"), nil
	case printFmtHeads:
		return []byte("heads"), nil
	default:
		return nil, fmt.Errorf("invalid format: %v", f)
	}
}

func printResults(s state, msg string, enumerations *big.Int) {
	switch s.printFmt {
	case printFmtNone:
	case printFmtShort:
		printShort(s, msg, enumerations)
	case printFmtLong:
		printLong(s, msg, enumerations)
	case printFmtSpace:
		printSpace(s, msg, enumerations)
	case printFmtHeads:
		printHeads(s, msg, enumerations)
	default:
		panic(fmt.Errorf("Invalid result format: %v", s.printFmt))
	}
}

func printShort(s state, msg string, enumerations *big.Int) {
	appendLine := func(strs []string, start, idx int, last byte) []string {
		diff := idx - start
		if diff > 0 && diff < 4 {
			for i := 0; i < diff; i++ {
				strs = append(strs, string([]byte{last}))
			}
		} else if diff >= 4 {
			strs = append(strs, fmt.Sprintf("%c(%d-%d)", last, start+1, idx))
		}
		return strs
	}

	var (
		strs  []string
		last  byte
		start int
	)
	for i, c := range s.coins {
		if c == last {
			continue
		}

		strs = appendLine(strs, start, i, last)
		last = c
		start = i
	}
	strs = appendLine(strs, start, len(s.coins), last)

	fmt.Printf("%s \t%s (%s)\n", strings.Join(strs, " "), msg, enumerations)
}

func printLong(s state, msg string, enumerations *big.Int) {
	fmt.Printf("%s %s (%s)\n", string(s.coins), msg, enumerations)
}

func printSpace(s state, msg string, enumerations *big.Int) {
	str := bytes.TrimLeft(bytes.Replace(s.coins, nil, []byte{' '}, -1), " ")
	fmt.Printf("%s %s (%s)\n", str, msg, enumerations)
}

func printHeads(s state, msg string, enumerations *big.Int) {
	maxHeads := s.heads * 2
	width := int(math.Log10(float64(len(s.coins)))) + 1
	idxFmt := fmt.Sprintf("%%-%dd", width)
	resFmt := fmt.Sprintf("%%-%ds", (width*maxHeads)+(maxHeads-1))

	idxs := make([]string, 0, maxHeads)
	for i, c := range s.coins {
		if c == 'H' {
			idxs = append(idxs, fmt.Sprintf(idxFmt, i+1))
		}
	}
	str := strings.Join(idxs, " ")
	fmt.Printf(resFmt+" %s (%s)\n", str, msg, enumerations)
}

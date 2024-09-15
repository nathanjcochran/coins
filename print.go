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
	printFmtLong
	printFmtHeads
)

func (f *printFmt) UnmarshalText(b []byte) error {
	s := string(b)
	switch s {
	case "none":
		*f = printFmtNone
	case "long":
		*f = printFmtLong
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
	case printFmtLong:
		return []byte("long"), nil
	case printFmtHeads:
		return []byte("heads"), nil
	default:
		return nil, fmt.Errorf("invalid format: %v", f)
	}
}

func printResults(s state, msg string, enumerations *big.Int) {
	switch s.printFmt {
	case printFmtNone:
	case printFmtLong:
		printLong(s, msg, enumerations)
	case printFmtHeads:
		printHeads(s, msg, enumerations)
	default:
		panic(fmt.Errorf("Invalid result format: %v", s.printFmt))
	}
}

func printLong(s state, msg string, enumerations *big.Int) {
	str := bytes.Replace(s.coins, nil, []byte{' '}, -1)
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

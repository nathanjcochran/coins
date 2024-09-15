package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/big"
	"slices"
)

var coinSides = []byte{'H', 'T'}

func main() {
	var (
		numCoins     = flag.Int("n", 10, "Number of coins to flip")
		printResults = flag.Bool("p", false, "Print results of each flip")
	)
	flag.Parse()

	draws, aliceWins, bobWins := flip(state{
		coins:        slices.Repeat([]byte{'_'}, *numCoins),
		draws:        big.NewInt(0),
		aliceWins:    big.NewInt(0),
		bobWins:      big.NewInt(0),
		printResults: *printResults,
	})
	fmt.Printf("Draws: %s\n", draws)
	fmt.Printf("Alice wins: %s\n", aliceWins)
	fmt.Printf("Bob wins: %s\n", bobWins)
}

type state struct {
	coins        []byte
	turn         int
	flipped      int
	bobIdx       int
	aliceIdx     int
	aliceHeads   int
	bobHeads     int
	draws        *big.Int
	aliceWins    *big.Int
	bobWins      *big.Int
	printResults bool
}

func flip(s state) (*big.Int, *big.Int, *big.Int) {
	// Whether Alice/Bob are actually flipping a new coin, or a coin that has
	// already been flipped by the other player
	aliceFlip := s.coins[s.aliceIdx] == '_'
	bobFlip := s.coins[s.bobIdx] == '_'

	switch {
	case aliceFlip && bobFlip && s.aliceIdx != s.bobIdx:
		s.flipped += 2
		for _, aliceCoin := range coinSides {
			s.coins[s.aliceIdx] = aliceCoin
			for _, bobCoin := range coinSides {
				s.coins[s.bobIdx] = bobCoin
				s.draws, s.aliceWins, s.bobWins = calculateResults(s)
			}
		}
	case aliceFlip:
		s.flipped += 1
		for _, aliceCoin := range coinSides {
			s.coins[s.aliceIdx] = aliceCoin
			s.draws, s.aliceWins, s.bobWins = calculateResults(s)
		}
	case bobFlip:
		s.flipped += 1
		for _, bobCoin := range coinSides {
			s.coins[s.bobIdx] = bobCoin
			s.draws, s.aliceWins, s.bobWins = calculateResults(s)
		}
	default:
		s.draws, s.aliceWins, s.bobWins = calculateResults(s)
	}

	// Reset coins that were flipped this round
	if aliceFlip {
		s.coins[s.aliceIdx] = '_'
	}
	if bobFlip {
		s.coins[s.bobIdx] = '_'
	}

	// Return accumulated win counts
	return s.draws, s.aliceWins, s.bobWins
}

func calculateResults(s state) (*big.Int, *big.Int, *big.Int) {
	if s.coins[s.aliceIdx] == 'H' {
		s.aliceHeads += 1
	}
	if s.coins[s.bobIdx] == 'H' {
		s.bobHeads += 1
	}

	switch {
	case s.aliceHeads == 2 && s.bobHeads == 2:
		// Entire subtree is a draw
		enumerations := remainingEnumerations(s)
		s.draws.Add(s.draws, enumerations)
		printResults(s, "Draw", enumerations)
		return s.draws, s.aliceWins, s.bobWins
	case s.aliceHeads == 2:
		// Alice wins entire subtree
		enumerations := remainingEnumerations(s)
		s.aliceWins.Add(s.aliceWins, enumerations)
		printResults(s, "Alice", enumerations)
		return s.draws, s.aliceWins, s.bobWins
	case s.bobHeads == 2:
		// Bob wins entire subtree
		enumerations := remainingEnumerations(s)
		s.bobWins.Add(s.bobWins, enumerations)
		printResults(s, "Bob", enumerations)
		return s.draws, s.aliceWins, s.bobWins
	case s.turn == len(s.coins)-1:
		enumerations := remainingEnumerations(s)
		s.draws.Add(s.draws, enumerations)
		printResults(s, "Draw", enumerations)
		return s.draws, s.aliceWins, s.bobWins
	default:
		// Still tied, flip another coin
		s.turn += 1
		s.aliceIdx += 1
		s.bobIdx = s.bobIdx + 2
		if s.bobIdx > len(s.coins)-1 {
			s.bobIdx = s.bobIdx - (len(s.coins) - 1)
		}
		return flip(s)
	}
}

func remainingEnumerations(s state) *big.Int {
	unFlipped := len(s.coins) - s.flipped

	return big.NewInt(0).Exp(
		big.NewInt(int64(len(coinSides))),
		big.NewInt(int64(unFlipped)),
		nil,
	)
}

func printResults(s state, msg string, enumerations *big.Int) {
	out := bytes.Replace(s.coins, nil, []byte{' '}, -1)
	fmt.Printf("%s %s (%s)\n", out, msg, enumerations)
}

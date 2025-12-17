// Package evaluator provides poker hand evaluation logic.
package evaluator

import "github.com/KyleKDang/poker-odds-engine/internal/card"

// Hand rank represents the strength of a poker hand.
// Using iota for auto-incrementing enum values.
type HandRank int

const (
	HighCard HandRank = iota + 1
	OnePair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

// HandRankNames maps ranks to their display names.
var HandRankNames = map[HandRank]string{
	HighCard:      "High Card",
	OnePair:       "One Pair",
	TwoPair:       "Two Pair",
	ThreeOfAKind:  "Three of a Kind",
	Straight:      "Straight",
	Flush:         "Flush",
	FullHouse:     "Full House",
	FourOfAKind:   "Four of a Kind",
	StraightFlush: "Straight Flush",
	RoyalFlush:    "Royal Flush",
}

// HandResult contains the evaluation result of a poker hand.
type HandResult struct {
	Rank    HandRank
	Label   string
	Kickers []int
}

// Compare compares two hand results.
// Returns: 1 if h1 wins, -1 if h2 wins, 0 if tie.
func (h1 *HandResult) Compare (h2 *HandResult) int {
	if h1.Rank != h2.Rank {
		if h1.Rank > h2.Rank {
			return 1
		}
		return -1
	}

	// Same rank - compare kickers
	for i := 0; i < len(h1.Kickers) && i < len(h2.Kickers); i++ {
		if h1.Kickers[i] > h2.Kickers[i] {
			return 1
		}
		if h1.Kickers[i] < h2.Kickers[i] {
			return -1
		}
	}

	return 0
}

// rankCounts counts how many of each rank appear in the hand.
func rankCounts(cards []*card.Card) map[card.Rank]int {
	counts := make(map[card.Rank]int)
	for _, c := range cards {
		counts[c.Rank]++
	}
	return counts
}

// isFlush checks if all cards have the same suit.
func isFlush(cards []*card.Card) bool {
	if len(cards) < 5 {
		return false
	}
	suit := cards[0].Suit
	for _, c := range cards[1:] {
		if c.Suit != suit {
			return false
		}
	}
	return true
}

// isStraight checks if cards form a straight.
// Returns whether it's a straight and the high card rank value.
func isStraight(cards []*card.Card) (bool, int) {
	if len(cards) < 5 {
		return false, 0
	}

	// Get unique rank values
	rankValues := make(map[int]bool)
	for _, c := range cards {
		rankValues[c.RankValue()] = true
	}

	// Convert to sorted slice (descending)
	values := make([]int, 0, len(rankValues))
	for v := range rankValues {
		values = append(values, v)
	}

	// Simple bubble sort (descending)
	for i := 0; i < len(values); i++ {
		for j := i + 1; j < len(values); j++ {
			if values[j] > values[i] {
				values[i], values[j] = values[j], values[i]
			}
		}
	}

	// Check for 5 consecutive values
	for i := 0; i <= len(values)-5; i++ {
		if values[i]-values[i+4] == 4 {
			return true, values[i]
		}
	}

	// Check for ace-low straight (A-2-3-4-5)
	if rankValues[12] && rankValues[0] && rankValues[1] && rankValues[2] && rankValues[3] {
		return true, 3
	}

	return false, 0
}

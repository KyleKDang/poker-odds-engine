package evaluator

import "github.com/KyleKDang/poker-odds-engine/internal/card"

// EvaluateHand finds the best 5-card poker hand from 1-7 cards.
func EvaluateHand(cards []*card.Card) *HandResult {
	if len(cards) < 1 {
		return nil
	}

	if len(cards) < 5 {
		return evaluateFiveCardHand(cards)
	}

	var bestHand *HandResult
	combinations := generateCombinations(cards, 5)
	
	for _, combo := range combinations {
		result := evaluateFiveCardHand(combo)
		if bestHand == nil || result.Compare(bestHand) > 0 {
			bestHand = result
		}
	}

	return bestHand
}

// evaluateFiveCardHand evaluates exactly 5 cards (or fewer for partial hands).
func evaluateFiveCardHand(cards []*card.Card) *HandResult {
	// Sort cards by rank value (highest first)
	sortedCards := make([]*card.Card, len(cards))
	copy(sortedCards, cards)
	
	for i := 0; i < len(sortedCards); i++ {
		for j := i + 1; j < len(sortedCards); j++ {
			if sortedCards[j].RankValue() > sortedCards[i].RankValue() {
				sortedCards[i], sortedCards[j] = sortedCards[j], sortedCards[i]
			}
		}
	}

	counts := rankCounts(sortedCards)
	flush := isFlush(sortedCards)
	straight, straightHigh := isStraight(sortedCards)

	// Build sorted list of rank counts for pattern matching
	type rankCount struct {
		rank  card.Rank
		count int
		value int
	}

	countsList := make([]rankCount, 0, len(counts))
	for rank, count := range counts {
		value := 0
		for i, r := range card.RankOrder {
			if r == rank {
				value = i
				break
			}
		}
		countsList = append(countsList, rankCount{rank, count, value})
	}

	// Sort by count (descending), then by rank value (descending)
	for i := 0; i < len(countsList); i++ {
		for j := i + 1; j < len(countsList); j++ {
			if countsList[j].count > countsList[i].count ||
				(countsList[j].count == countsList[i].count && countsList[j].value > countsList[i].value) {
				countsList[i], countsList[j] = countsList[j], countsList[i]
			}
		}
	}

	// Check for each hand type (best to worst)
	
	// Royal Flush
	if straight && flush && straightHigh == 12 {
		return &HandResult{
			Rank:    RoyalFlush,
			Label:   HandRankNames[RoyalFlush],
			Kickers: []int{},
		}
	}

	// Straight Flush
	if straight && flush {
		return &HandResult{
			Rank:    StraightFlush,
			Label:   HandRankNames[StraightFlush],
			Kickers: []int{straightHigh},
		}
	}

	// Four of a Kind
	if len(sortedCards) >= 4 && countsList[0].count == 4 {
		kickers := []int{countsList[0].value}
		if len(countsList) > 1 {
			kickers = append(kickers, countsList[1].value)
		}
		return &HandResult{
			Rank:    FourOfAKind,
			Label:   HandRankNames[FourOfAKind],
			Kickers: kickers,
		}
	}

	// Full House
	if len(sortedCards) >= 5 && countsList[0].count == 3 && countsList[1].count >= 2 {
		return &HandResult{
			Rank:    FullHouse,
			Label:   HandRankNames[FullHouse],
			Kickers: []int{countsList[0].value, countsList[1].value},
		}
	}

	// Flush
	if flush {
		kickers := make([]int, 0, len(sortedCards))
		for _, c := range sortedCards {
			kickers = append(kickers, c.RankValue())
		}
		return &HandResult{
			Rank:    Flush,
			Label:   HandRankNames[Flush],
			Kickers: kickers,
		}
	}

	// Straight
	if straight {
		return &HandResult{
			Rank:    Straight,
			Label:   HandRankNames[Straight],
			Kickers: []int{straightHigh},
		}
	}

	// Three of a Kind
	if len(sortedCards) >= 3 && countsList[0].count == 3 {
		kickers := []int{countsList[0].value}
		for i := 1; i < len(countsList); i++ {
			kickers = append(kickers, countsList[i].value)
		}
		return &HandResult{
			Rank:    ThreeOfAKind,
			Label:   HandRankNames[ThreeOfAKind],
			Kickers: kickers,
		}
	}

	// Two Pair
	if len(sortedCards) >= 4 && countsList[0].count == 2 && countsList[1].count == 2 {
		kickers := []int{countsList[0].value, countsList[1].value}
		if len(countsList) > 2 {
			kickers = append(kickers, countsList[2].value)
		}
		return &HandResult{
			Rank:    TwoPair,
			Label:   HandRankNames[TwoPair],
			Kickers: kickers,
		}
	}

	// One Pair
	if len(sortedCards) >= 2 && countsList[0].count == 2 {
		kickers := []int{countsList[0].value}
		for i := 1; i < len(countsList); i++ {
			kickers = append(kickers, countsList[i].value)
		}
		return &HandResult{
			Rank:    OnePair,
			Label:   HandRankNames[OnePair],
			Kickers: kickers,
		}
	}

	// High Card
	kickers := make([]int, 0, len(sortedCards))
	for _, c := range sortedCards {
		kickers = append(kickers, c.RankValue())
	}
	return &HandResult{
		Rank:    HighCard,
		Label:   HandRankNames[HighCard],
		Kickers: kickers,
	}
}

// generateCombinations generates all k-size combinations from cards.
func generateCombinations(cards []*card.Card, k int) [][]*card.Card {
	var result [][]*card.Card
	n := len(cards)

	// Recursive helper function
	var helper func(start int, combo []*card.Card)
	helper = func(start int, combo []*card.Card) {
		if len(combo) == k {
			comb := make([]*card.Card, k)
			copy(comb, combo)
			result = append(result, comb)
			return
		}

		for i := start; i < n; i++ {
			helper(i+1, append(combo, cards[i]))
		}
	}

	helper(0, []*card.Card{})
	return result
}

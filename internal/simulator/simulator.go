// Package simulator provides Monte Carlo poker odds calculation.
package simulator

import (
	"math/rand"
	"sync"
	"time"

	"github.com/KyleKDang/poker-odds-engine/internal/card"
	"github.com/KyleKDang/poker-odds-engine/internal/evaluator"
)

// OddsResult contains win/tie/loss probabilities.
type OddsResult struct {
	Win  float64 `json:"win"`
	Tie  float64 `json:"tie"`
	Loss float64 `json:"loss"`
}

// CalculateOdds runs Monte Carlo simulation to calculate poker odds.
func CalculateOdds(holeCards, boardCards []*card.Card, numOpponents, simulations, workers int) *OddsResult {
	if workers < 1 {
		workers = 4
	}
	if simulations < 1 {
		simulations = 10000
	}

	simulationsPerWorker := simulations / workers
	extraSims := simulations % workers

	var wg sync.WaitGroup
	results := make(chan workerResult, workers)

	// Launch worker goroutines
	for i := 0; i < workers; i++ {
		wg.Add(1)

		sims := simulationsPerWorker
		if i < extraSims {
			sims++
		}

		go func() {
			defer wg.Done()
			result := runSimulations(holeCards, boardCards, numOpponents, sims)
			results <- result
		}()
	}

	// Close channel when all workers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Aggregate results
	totalWins := 0
	totalTies := 0
	totalSims := 0

	for result := range results {
		totalWins += result.wins
		totalTies += result.ties
		totalSims += result.simulations
	}

	totalLosses := totalSims - totalWins - totalTies

	return &OddsResult{
		Win: float64(totalWins) / float64(totalSims),
		Tie: float64(totalTies) / float64(totalSims),
		Loss: float64(totalLosses) / float64(totalSims),
	}
}

// workerResult holds results from a single worker goroutine.
type workerResult struct {
	wins        int
	ties        int
	simulations int
}

// runSimulations performs Monte Carlo simulations for one worker.
func runSimulations(holeCards, boardCards []*card.Card, numOpponents, simulations int) workerResult {
	known := append(holeCards, boardCards...)
	deck := card.RemoveCards(card.NewDeck(), known)

	wins := 0
	ties := 0

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Run simulations
	for i := 0; i < simulations; i++ {
		shuffleDeck(deck, rng)

		missingCards := 5 - len(boardCards)
		fullBoard := make([]*card.Card, len(boardCards))
		copy(fullBoard, boardCards)
		fullBoard = append(fullBoard, deck[:missingCards]...)

		opponentHands := make([][]*card.Card, numOpponents)
		idx := missingCards
		for j := 0; j < numOpponents; j++ {
			opponentHands[j] = []*card.Card{deck[idx], deck[idx+1]}
			idx += 2
		}

		playerCards := append(holeCards, fullBoard...)
		playerResult := evaluator.EvaluateHand(playerCards)

		var bestOpponent *evaluator.HandResult
		for _, oppHole := range opponentHands {
			oppCards := append(oppHole, fullBoard...)
			oppResult := evaluator.EvaluateHand(oppCards)

			if bestOpponent == nil || oppResult.Compare(bestOpponent) > 0 {
				bestOpponent = oppResult
			}
		}

		comparison := playerResult.Compare(bestOpponent)
		if comparison > 0 {
			wins++
		} else if comparison == 0 {
			ties++
		}
	}

	return workerResult{
		wins:        wins,
		ties:        ties,
		simulations: simulations,
	}
}

// shuffleDeck shuffles a deck in place using Fisher-Yates algorithm.
func shuffleDeck(deck []*card.Card, rng *rand.Rand) {
	for i := len(deck) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		deck[i], deck[j] = deck[j], deck[i]
	}
}

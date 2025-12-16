// Package card provides playing card operations for poker.
package card

import (
	"fmt"
	"strings"
)

// Rank represents a card rank (2-A).
type Rank string

// Suit represents a card suit (S/H/D/C).
type Suit string

const (
	Two   Rank = "2"
	Three Rank = "3"
	Four  Rank = "4"
	Five  Rank = "5"
	Six   Rank = "6"
	Seven Rank = "7"
	Eight Rank = "8"
	Nine  Rank = "9"
	Ten   Rank = "T"
	Jack  Rank = "J"
	Queen Rank = "Q"
	King  Rank = "K"
	Ace   Rank = "A"
)

const (
	Spades   Suit = "S"
	Hearts   Suit = "H"
	Diamonds Suit = "D"
	Clubs    Suit = "C"
)

// RankOrder defines ranks from lowest to highest.
var RankOrder = []Rank{
	Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King, Ace,
}

// AllSuits contains all four suits.
var AllSuits = []Suit{Spades, Hearts, Diamonds, Clubs}

// Card represents a playing card.
type Card struct {
	Rank Rank
	Suit Suit
}

// NewCard creates a card from a 2-character code (e.g., "AS").
func NewCard(code string) (*Card, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	if len(code) != 2 {
		return nil, fmt.Errorf("invalid card code: %s", code)
	}

	rank := Rank(code[0:1])
	suit := Suit(code[1:2])

	validRank := false
	for _, r := range RankOrder {
		if r == rank {
			validRank = true
			break
		}
	}
	if !validRank {
		return nil, fmt.Errorf("invalid rank: %s", rank)
	}

	validSuit := false
	for _, s := range AllSuits {
		if s == suit {
			validSuit = true
			break
		}
	}
	if !validSuit {
		return nil, fmt.Errorf("invalid suit: %s", suit)
	}

	return &Card{Rank: rank, Suit: suit}, nil
}

// String returns the card's string representation.
func (c *Card) String() string {
	return string(c.Rank) + string(c.Suit)
}

// RankValue returns the numeric rank value (0-12).
func (c *Card) RankValue() int {
	for i, r := range RankOrder {
		if r == c.Rank {
			return i
		}
	}
	return -1
}

// Equal checks if two cards are identical.
func (c *Card) Equal(other *Card) bool {
	return c.Rank == other.Rank && c.Suit == other.Suit
}

// NewDeck creates a standard 52-card deck.
func NewDeck() []*Card {
	deck := make([]*Card, 0, 52)
	for _, suit := range AllSuits {
		for _, rank := range RankOrder {
			deck = append(deck, &Card{Rank: rank, Suit: suit})
		}
	}
	return deck
}

// RemoveCards returns a deck with specified cards removed.
func RemoveCards(deck []*Card, toRemove []*Card) []*Card {
	result := make([]*Card, 0, len(deck))
	for _, card := range deck {
		keep := true
		for _, remove := range toRemove {
			if card.Equal(remove) {
				keep = false
				break
			}
		}
		if keep {
			result = append(result, card)
		}
	}
	return result
}

// ParseCards converts string codes to Card objects.
func ParseCards(codes []string) ([]*Card, error) {
	cards := make([]*Card, 0, len(codes))
	for _, code := range codes {
		card, err := NewCard(code)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, nil
}

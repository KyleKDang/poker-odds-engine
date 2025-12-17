// Package models defines API request and response structures.
package models

// EvaluateRequest contains cards to evaluate.
type EvaluateRequest struct {
	HoleCards  []string `json:"hole_cards" binding:"required"`
	BoardCards []string `json:"board_cards" binding:"required"`
}

// EvaluateResponse contains the evaluated hand result.
type EvaluateResponse struct {
	Hand string `json:"hand"`
	Rank int    `json:"rank"`
}

// OddsRequest contains parameters for odds calculation.
type OddsRequest struct {
	HoleCards    []string `json:"hole_cards" binding:"required"`
	BoardCards   []string `json:"board_cards" binding:"required"`
	NumOpponents int      `json:"num_opponents" binding:"required,min=1,max=9"`
	Simulations  int      `json:"simulations,omitempty"`
	Workers      int      `json:"workers,omitempty"`
}

// OddsResponse contains calculated odds.
type OddsResponse struct {
	Win  float64 `json:"win"`
	Tie  float64 `json:"tie"`
	Loss float64 `json:"loss"`
}

// ErrorResponse contains error information.
type ErrorResponse struct {
	Error string `json:"error"`
}

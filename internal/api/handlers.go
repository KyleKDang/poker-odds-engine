// Package api provides HTTP handlers for the poker odds engine.
package api

import (
	"net/http"

	"github.com/KyleKDang/poker-odds-engine/internal/card"
	"github.com/KyleKDang/poker-odds-engine/internal/evaluator"
	"github.com/KyleKDang/poker-odds-engine/internal/simulator"
	"github.com/KyleKDang/poker-odds-engine/pkg/models"
	"github.com/gin-gonic/gin"
)

// HandleHealth returns server health status.
func HandleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"service": "poker-odds-engine",
	})
}

// HandleEvaluate evaluates a poker hand.
func HandleEvaluate(c *gin.Context) {
	var req models.EvaluateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request: " + err.Error(),
		})
		return
	}

	holeCards, err := card.ParseCards(req.HoleCards)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid hole cards: " + err.Error(),
		})
		return
	}

	boardCards, err := card.ParseCards(req.BoardCards)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid board cards: " + err.Error(),
		})
		return
	}

	allCards := append(holeCards, boardCards...)
	result := evaluator.EvaluateHand(allCards)

	if result == nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Unable to evaluate hand",
		})
		return
	}

	c.JSON(http.StatusOK, models.EvaluateResponse{
		Hand: result.Label,
		Rank: int(result.Rank),
	})
}

// HandleOdds calculates winning odds using Monte Carlo simulation.
func HandleOdds(c *gin.Context) {
	var req models.OddsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid request: " + err.Error(),
		})
		return
	}

	if req.Simulations <= 0 {
		req.Simulations = 10000
	}
	if req.Workers <= 0 {
		req.Workers = 4
	}

	holeCards, err := card.ParseCards(req.HoleCards)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid hole cards: " + err.Error(),
		})
		return
	}

	boardCards, err := card.ParseCards(req.BoardCards)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid board cards: " + err.Error(),
		})
		return
	}

	if len(holeCards) != 2 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Must provide exactly 2 hole cards",
		})
		return
	}
	if len(boardCards) > 5 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Board cannot have more than 5 cards",
		})
		return
	}
	
	result := simulator.CalculateOdds(holeCards, boardCards, req.NumOpponents, req.Simulations, req.Workers)

	c.JSON(http.StatusOK, models.OddsResponse{
		Win:  result.Win,
		Tie:  result.Tie,
		Loss: result.Loss,
	})
}

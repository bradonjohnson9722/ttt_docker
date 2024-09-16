package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Start a new game
func StartGame(c *gin.Context) {
	// Create a new game with a unique ID
	game := Game{
		ID:     uuid.New().String(),
		Board:  [9]string{"", "", "", "", "", "", "", "", ""},
		Status: "ongoing",
	}

	// Save the new game to the database
	err := saveGame(game)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create game"})
		return
	}

	// Send back the initial game state
	c.JSON(200, game)
}

// Make a move (called by the frontend after the player makes a move)
func MakeMove(c *gin.Context) {
	var move Move
	var game Game

	// Decode incoming move from JSON
	if err := c.BindJSON(&move); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Get game ID from query params (you may also pass it via body if preferred)
	gameID := c.Query("id")

	// Find the game in the database
	game, err := findGameByID(gameID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Game not found"})
		return
	}

	// Process the move (update board, check win conditions, and let the bot play)
	updatedBoard, status := playTurn(game.Board, move.Position)
	game.Board = updatedBoard
	game.Status = status

	// Update the game in the database
	err = updateGame(game)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update game"})
		return
	}

	// Return the updated game state
	c.JSON(200, game)
}

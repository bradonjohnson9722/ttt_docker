package main

import (
	"math/rand"
	"time"
)

// Function to check for a win condition
func checkWin(board [9]string, player string) bool {
	winPatterns := [8][3]int{
		{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // Rows
		{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // Columns
		{0, 4, 8}, {2, 4, 6}, // Diagonals
	}

	for _, pattern := range winPatterns {
		if board[pattern[0]] == player && board[pattern[1]] == player && board[pattern[2]] == player {
			return true
		}
	}

	return false
}

// Bot logic to pick a random empty spot
func getBotMove(board [9]string) int {
	emptyPositions := []int{}

	// Find empty positions on the board
	for i, spot := range board {
		if spot == "" {
			emptyPositions = append(emptyPositions, i)
		}
	}

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Pick a random empty position
	if len(emptyPositions) > 0 {
		return emptyPositions[rand.Intn(len(emptyPositions))]
	}

	// No empty positions left, return -1
	return -1
}

// Function to process a turn (player and bot)
func playTurn(board [9]string, playerMove int) ([9]string, string) {
	// Player makes a move
	board[playerMove] = "X"

	// Check if the player won
	if checkWin(board, "X") {
		return board, "Player wins"
	}

	// Bot makes a move
	botMove := getBotMove(board)
	if botMove != -1 {
		board[botMove] = "O"

		// Check if the bot won
		if checkWin(board, "O") {
			return board, "Bot wins"
		}
	} else {
		return board, "Draw"
	}

	return board, "ongoing"
}

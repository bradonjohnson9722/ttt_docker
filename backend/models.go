package main

// Move represents a move made by a player or the bot
type Move struct {
	Player   string `json:"player"`   // "user" or "bot"
	Position int    `json:"position"` // Position of the move (0-8 for the grid)
}

// Game represents a Tic-Tac-Toe game
type Game struct {
	ID     string    `json:"id" bson:"id"`         // Unique game ID
	Board  [9]string `json:"board" bson:"board"`   // 3x3 Tic-Tac-Toe board
	Status string    `json:"status" bson:"status"` // "ongoing", "won", "draw"
	Flag   bool
}

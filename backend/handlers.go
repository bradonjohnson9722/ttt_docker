package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

var rdb *redis.Client

// Function to sed message t RabbitMQ
func sendEndGameMessage() error {
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://admin:admin@rabbitmq:5672")
	fmt.Println("66666666666666666666666666666666", err)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declare the queue
	queue, err := ch.QueueDeclare(
		"game_events", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return err
	}

	// Create the message
	message := "end_brandon_game"

	// Publish the message
	err = ch.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Sent: %s", message)
	http.Get("http://192.168.1.151:1978/read_message")
	return nil
}

// Start a new game
func StartGame(c *gin.Context) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "192.168.1.151:1988", // Redis server address
		Password: "",                   // Set password if required
		DB:       0,                    // Use default DB
	})
	userID := c.Query("userid")
	tokenJSON, err := rdb.Get(ctx, userID).Result()
	fmt.Println("333333333333333333", userID)
	fmt.Println("444444444444444444", tokenJSON)

	if err != nil {
		return
	}
	// Create a new game with a unique ID
	if tokenJSON != "" {
		game := Game{
			ID:     uuid.New().String(),
			Board:  [9]string{"", "", "", "", "", "", "", "", ""},
			Status: "ongoing",
			Flag:   true,
		}
		err := saveGame(game)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to create game"})
			return
		}
		c.JSON(200, game)

	} else {
		c.JSON(500, gin.H{"error": "Failed to create game"})
	}

	// Save the new game to the database

	// Send back the initial game state

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
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: "192.168.1.151:1988", // Change to your Redis server address
	})
	exists, err := client.Exists(ctx, "104297523855792674593").Result()
	if exists == 1 {
		game.Flag = true
	} else {
		game.Flag = false
	}
	// Update the game in the database
	err = updateGame(game)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update game"})
		return
	}

	// if the game is finished, the auth must logout
	if status != "ongoing" {
		// Simulating the end of a game
		err := sendEndGameMessage()
		if err != nil {
			log.Fatalf("Failed to send message: %s", err)
		}
	}

	// Return the updated game state
	c.JSON(200, game)

}

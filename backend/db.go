package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var gameCollection *mongo.Collection

// Initialize MongoDB connection
func InitDB() {
	clientOptions := options.Client().ApplyURI("mongodb://root:password@mongo:27017") // MongoDB URI
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Select the database and collection
	gameCollection = client.Database("tic-tac-toe").Collection("games")
}

// Save a new game to the database
func saveGame(game Game) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := gameCollection.InsertOne(ctx, game)
	return err
}

// Find game by ID from the database
func findGameByID(id string) (Game, error) {
	var game Game
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err := gameCollection.FindOne(ctx, bson.M{"id": id}).Decode(&game)
	return game, err
}

// Update an existing game in the database
func updateGame(game Game) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	filter := bson.M{"id": game.ID}
	_, err := gameCollection.ReplaceOne(ctx, filter, game)
	return err
}

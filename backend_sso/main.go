package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:1978/auth/google/callback",
		ClientID:     "280247170993-pj0gv5dmpj8l6cukvdokisdgel4diva0.apps.googleusercontent.com",
		ClientSecret: "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",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	randomState = "random" // You should generate a random state for security
)

func main() {

	// Initialize Gin router
	router := gin.Default()

	// Setup CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Frontend origin
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	router.GET("/auth/google//login", GoogleLogin)
	router.GET("/auth/google/callback", GoogleCallback)
	router.GET("/read_message", Read_RabbitMQ_message)

	// Start the server
	log.Println("Backend-Server-sso starting on port 1978...")
	router.Run(":1978") // Server will run on port 1978
}

func GoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL(randomState)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	ctx := context.Background()

	if c.Query("state") != randomState {
		c.JSON(http.StatusBadRequest, "Invalid state")
		return
	}
	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)

	// Create a new Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.1.151:1988", // Redis server address
		Password: "",                   // No password set
		DB:       0,                    // Use default DB
	})

	// end create redis client
	if err != nil {
		c.JSON(http.StatusBadRequest, "Couldn't get token")
		return
	}

	client := googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusBadRequest, "Couldn't get user info")
		return
	}
	// defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// ex_body := body
	bodyJSON := string(body)
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(bodyJSON), &result); err != nil {
		fmt.Errorf("couldn't decode JSON: %w", err)
	}

	// Push token to Redis
	tokenJSON, err := json.Marshal(token)
	if err != nil {
		return
	}
	err = rdb.Set(ctx, result["id"].(string), tokenJSON, 3600*time.Second).Err() // 1 hour expiration
	if err != nil {
		log.Fatalf("Could not set token: %v", err)
	}

	// Here, you can save the user info to your MongoDB database
	// and create a session if necessary
	c.JSON(http.StatusOK, resp)
	// c.Redirect(http.StatusMovedPermanently, "http://192.168.1.151:1975?"+string(ex_body))
	c.Redirect(http.StatusMovedPermanently, "http://192.168.1.151:1975?"+result["id"].(string)+result["name"].(string))
}

func Read_RabbitMQ_message(c *gin.Context) {
	// Connect to RabbitMQ server
	conn, err := amqp.Dial("amqp://admin:admin@rabbitmq:5672")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
	}
	defer channel.Close()

	queueName := "game_events"

	// Declare the queue
	_, err = channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	// Set up a consumer
	msgs, err := channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	// Start receiving messages
	go func() {
		for msg := range msgs {
			fmt.Printf("Received message777777777777777777777777: %s\n", msg.Body)
		}
	}()

	fmt.Println("Waiting for messages. To exit press CTRL+C")
	// Block forever
	// select {}

	// delete the key-value in redis
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: "192.168.1.151:1988", // Change to your Redis server address
	})

	// Delete a key
	err = client.Del(ctx, "104297523855792674593").Err()
	if err != nil {
		log.Fatalf("Could not delete key: %v", err)
	}

}

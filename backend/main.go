package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:1972/auth/google/callback",
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
	// Initialize MongoDB connection
	InitDB()

	// Initialize Gin router
	router := gin.Default()

	// Setup CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Frontend origin
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	// Define API endpoints
	router.POST("/start-game", StartGame)
	router.POST("/make-move", MakeMove)

	router.GET("/auth/google//login", GoogleLogin)
	router.GET("/auth/google/callback", GoogleCallback)

	// Start the server
	log.Println("Server starting on port 1972...")
	router.Run(":1972") // Server will run on port 8080
}

func GoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL(randomState)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	if c.Query("state") != randomState {
		c.JSON(http.StatusBadRequest, "Invalid state")
		return
	}
	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	fmt.Print("222222222222222", err)
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
	fmt.Println("1111111111111111111111111111111", string(body))
	// Here, you can save the user info to your MongoDB database
	// and create a session if necessary
	c.JSON(http.StatusOK, resp)
	c.Redirect(http.StatusMovedPermanently, "http://192.168.1.151:1975?"+string(body))
}

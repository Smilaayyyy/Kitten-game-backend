package main

import (
	"fmt"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"backend/handlers"
	"backend/utils"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// Initialize Redis
	err = utils.InitRedis()
	if err != nil {
		fmt.Println("Error initializing Redis:", err)
		return
	}

	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true, 
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	// Define HTTP routes
	r.POST("/api/register", handlers.RegisterHandler)

	r.POST("/api/start-game", handlers.StartGame)
	r.POST("/api/draw-card", handlers.DrawCard)
	r.GET("/api/leaderboard", handlers.GetLeaderboard)

	
    
	// Additional routes
	r.POST("/api/save-game", handlers.SaveGameHandler)
	r.POST("/record-win", handlers.RecordWinHandler)
	r.GET("/leaderboard", handlers.LeaderboardHandler)

	
    r.GET("/ws/leaderboard", func(c *gin.Context) {
        handlers.LeaderboardWebSocketHandler(c)
    })

    // Start WebSocket broadcasting
    go handlers.BroadcastLeaderboardUpdates()

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server run failed:", err)
	}
}

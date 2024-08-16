package handlers

import (
    "backend/utils"
    "net/http"
    "github.com/gin-gonic/gin"
    "fmt"
    "log"
    
)
// RegisterHandler handles user registration
func RegisterHandler(c *gin.Context) {
    var req map[string]string
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    username, ok := req["username"]
    if !ok || username == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
        return
    }

    // Save username in Redis with initial wins = 0
    err := utils.Rdb.Set(utils.Ctx, username, 0, 0).Err()
    if err != nil {
        fmt.Printf("Redis Set error: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save username"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "registered"})
}

func RecordWinHandler(c *gin.Context) {
    var req map[string]interface{}
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
        return
    }

    username, ok := req["username"].(string)
    if !ok || username == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
        return
    }

    points, ok := req["points"].(float64)
    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Points are required and must be a number"})
        return
    }

    if points == 2 {
        // Increment the win count in Redis
        if err := utils.Rdb.Incr(utils.Ctx, username).Err(); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record win"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"status": "win recorded"})
    } else {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Win not recorded. Points must be equal to 2."})
    }
}


// LeaderboardHandler fetches the leaderboard
func LeaderboardHandler(c *gin.Context) {
    users, err := utils.Rdb.Keys(utils.Ctx, "*").Result()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaderboard"})
        return
    }

    leaderboard := make(map[string]int)
    for _, user := range users {
        wins, err := utils.Rdb.Get(utils.Ctx, user).Int()
        if err != nil {
            continue // Skip users with issues
        }
        leaderboard[user] = wins
    }

    // Log the leaderboard data
    log.Printf("Leaderboard data: %v", leaderboard)

    // Send the leaderboard data as an HTTP response
    c.JSON(http.StatusOK, gin.H{"leaderboard": leaderboard})

    // Broadcast the leaderboard data over WebSocket
    leaderboardBroadcast <- leaderboard
}

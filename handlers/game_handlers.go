package handlers

import (
    "net/http"
    "sync"
    "github.com/gin-gonic/gin"
    
)

// GameState represents the game state structure
type GameState struct {
    Username    string   `json:"username"`
    Deck        []string `json:"deck"`
    DrawnCard   string   `json:"drawnCard"`
    Points      int      `json:"points"`
    DefuseCount int      `json:"defuseCount"`
    GameOver    bool     `json:"gameOver"`
    Message     string   `json:"message"`
}

// In-memory database (simulated)
var gameStateDB = make(map[string]GameState)
var mu sync.Mutex

// SaveGameHandler handles the saving of the game state using Gin context
func SaveGameHandler(c *gin.Context) {
    var gameState GameState
    if err := c.ShouldBindJSON(&gameState); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    mu.Lock()
    gameStateDB[gameState.Username] = gameState
    mu.Unlock()

    c.String(http.StatusOK, "Game state saved successfully!")
}

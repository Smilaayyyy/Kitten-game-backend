package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"backend/models"
	"backend/utils"
)

// UpdateLeaderboard handles updating the leaderboard with the latest scores
func UpdateLeaderboard(username string, points int) {
	// Retrieve the existing user data
	userJson, err := utils.Rdb.Get(utils.Ctx, "user:"+username).Result()
	if err != nil {
		return // Handle error if needed
	}

	var user models.User
	err = json.Unmarshal([]byte(userJson), &user)
	if err != nil {
		return // Handle error if needed
	}

	// Update user points
	user.Points = points // Assuming User struct has a Points field

	// Save the updated user data
	updatedUserJson, err := json.Marshal(user)
	if err != nil {
		return // Handle error if needed
	}

	utils.Rdb.Set(utils.Ctx, "user:"+username, string(updatedUserJson), 0).Err()

	// Broadcast the updated leaderboard
	broadcastLeaderboard()
}

// broadcastLeaderboard broadcasts the current leaderboard to all WebSocket clients
func broadcastLeaderboard() {
	var leaderboard []models.User

	iter := utils.Rdb.Scan(utils.Ctx, 0, "user:*", 0).Iterator()
	for iter.Next(utils.Ctx) {
		userJson, err := utils.Rdb.Get(utils.Ctx, iter.Val()).Result()
		if err != nil {
			continue
		}

		var user models.User
		err = json.Unmarshal([]byte(userJson), &user)
		if err != nil {
			continue
		}

		leaderboard = append(leaderboard, user)
	}
	if err := iter.Err(); err != nil {
		return // Handle error if needed
	}

	// Convert leaderboard to JSON
	leaderboardJson, err := json.Marshal(leaderboard)
	if err != nil {
		return // Handle error if needed
	}

	// Broadcast to all connected WebSocket clients
	utils.WebSocketClients.Broadcast(leaderboardJson)
}

// GetLeaderboard retrieves and returns the current leaderboard
func GetLeaderboard(c *gin.Context) {
	var leaderboard []models.User

	iter := utils.Rdb.Scan(utils.Ctx, 0, "user:*", 0).Iterator()
	for iter.Next(utils.Ctx) {
		userJson, err := utils.Rdb.Get(utils.Ctx, iter.Val()).Result()
		if err != nil {
			continue
		}

		var user models.User
		err = json.Unmarshal([]byte(userJson), &user)
		if err != nil {
			continue
		}

		leaderboard = append(leaderboard, user)
	}
	if err := iter.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve leaderboard"})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}

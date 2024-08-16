package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
	"backend/utils"
	"github.com/go-redis/redis/v8"
	"strings"
)

type Card struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

func serializeDeck(deck []Card) string {
	var sb strings.Builder
	for _, card := range deck {
		sb.WriteString(fmt.Sprintf("%d:%s,", card.ID, card.Type))
	}
	return strings.TrimRight(sb.String(), ",")
}

func deserializeDeck(deckStr string) []Card {
	var deck []Card
	cardStrs := strings.Split(deckStr, ",")
	for _, cardStr := range cardStrs {
		var card Card
		fmt.Sscanf(cardStr, "%d:%s", &card.ID, &card.Type)
		deck = append(deck, card)
	}
	return deck
}

func StartGame(c *gin.Context) {
	deck := []Card{
		{ID: 1, Type: "Cat"},
		{ID: 2, Type: "Cat"},
		{ID: 3, Type: "Defuse"},
		{ID: 4, Type: "Shuffle"},
		{ID: 5, Type: "Exploding Kitten"},
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	deckStr := serializeDeck(deck)

	err := utils.Rdb.Set(utils.Ctx, "deck", deckStr, 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save deck data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deck": deck})
}
func DrawCard(c *gin.Context) {
	username := c.Param("username")

	// Retrieve deck from Redis
	deckStr, err := utils.Rdb.Get(utils.Ctx, "deck").Result()
	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No game in progress"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve deck data"})
		return
	}

	// Deserialize deck
	deck := deserializeDeck(deckStr)
	if len(deck) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Deck is empty"})
		return
	}

	// Draw the top card
	drawnCard := deck[0]
	deck = deck[1:]

	// Handle special cards
	if drawnCard.Type == "Exploding Kitten" {
		defuseCardStr, err := utils.Rdb.Get(utils.Ctx, "defuse:"+username).Result()
		if err == nil && defuseCardStr == "present" {
			// Player has a defuse card, allow to continue
			utils.Rdb.Del(utils.Ctx, "defuse:"+username) // Remove one defuse card
			c.JSON(http.StatusOK, gin.H{"card": drawnCard, "gameStatus": "defused"})
			return
		}
		utils.Rdb.Set(utils.Ctx, "gameStatus:"+username, "game_over", 0)
		c.JSON(http.StatusOK, gin.H{"card": drawnCard, "gameStatus": "game_over"})
		return
	}

	if drawnCard.Type == "Defuse" {
		utils.Rdb.Set(utils.Ctx, "defuse:"+username, "present", 0) // Store defuse card for player
	}

	if drawnCard.Type == "Shuffle" {
		// Restart the game (reshuffle the deck)
		StartGame(c)
		return
	}

	// Update the deck in Redis
	deckStr = serializeDeck(deck)
	err = utils.Rdb.Set(utils.Ctx, "deck", deckStr, 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save deck data"})
		return
	}

	// Check if the player has won
	if len(deck) == 0 {
		utils.Rdb.Set(utils.Ctx, "gameStatus:"+username, "won", 0)
		c.JSON(http.StatusOK, gin.H{"card": drawnCard, "gameStatus": "won"})
		return
	}

	// Game is still in progress
	c.JSON(http.StatusOK, gin.H{"card": drawnCard, "gameStatus": "in_progress"})
}

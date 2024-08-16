// handlers/websockethandler.go

package handlers

import (
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "log"
    "net/http"
    
)

var leaderboardClients = make(map[*websocket.Conn]bool)
var leaderboardBroadcast = make(chan map[string]int)

// Upgrade HTTP connection to WebSocket
var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func LeaderboardWebSocketHandler(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Println("Error upgrading to WebSocket:", err)
        return
    }
    defer conn.Close()

    leaderboardClients[conn] = true

    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            log.Println("Error reading WebSocket message:", err)
            delete(leaderboardClients, conn)
            break
        }
    }
}

func BroadcastLeaderboardUpdates() {
    for {
        leaderboard := <-leaderboardBroadcast
        for client := range leaderboardClients {
            err := client.WriteJSON(leaderboard)
            if err != nil {
                log.Println("Error writing JSON to WebSocket:", err)
                client.Close()
                delete(leaderboardClients, client)
            }
        }
    }
}

package models

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Points   int    `json:"points"`
}
type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}
package model

// Users For simplicity, we'll use an in-memory store. In a real application, you'd use a database.
var Users = make(map[string]string)

// User represents a user in our system
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

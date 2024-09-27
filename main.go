package main

import (
	auth "gin-img-processer/auth"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in our system
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// For simplicity, we'll use an in-memory store. In a real application, you'd use a database.
var users = make(map[string]string)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	public := r.Group("/public")
	{
		public.POST("/signup", signUp)
		public.GET("/signup", showSignupForm)
		public.POST("/login", logIn)

	}
	private := r.Group("/api")
	private.Use(auth.Authenticate())

	r.Run(":8080")
}
func showSignupForm(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", gin.H{})
}

func signUp(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	if _, exists := users[user.Username]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	// Store the user
	users[user.Username] = string(hashedPassword)

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func logIn(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the hashed password
	storedPassword, exists := users[user.Username]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Compare the stored hashed password with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully signed in"})
}

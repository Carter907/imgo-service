package main

import (
	"gin-img-processer/auth"
	"gin-img-processer/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

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
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	if _, exists := model.Users[user.Username]; exists {
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
	model.Users[user.Username] = string(hashedPassword)

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func logIn(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the hashed password
	storedPassword, exists := model.Users[user.Username]
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

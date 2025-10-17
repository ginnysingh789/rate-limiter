package Middleware

import (
	"net/http"
	models "rate-limiter/Models"

	"github.com/gin-gonic/gin"
)

func CheckUserInput(c *gin.Context) {
	var requestBody models.User

	if err := c.BindJSON(&requestBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}
	if requestBody.Username == "" && requestBody.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Please enter username and password"})
		return
	} else if requestBody.Username == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Please enter username"})
		return
	} else if requestBody.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Please enter password"})
		return
	}
	c.Set("validatedUser", requestBody)

	c.Next()
}

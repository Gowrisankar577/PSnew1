package handles

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func GetUserProfileImage(c *gin.Context) {
	userId := c.Query("userId")
	if userId == "" {
		c.JSON(400, gin.H{"error": "userId required"})
		return
	}

	imagePath := fmt.Sprintf("images/users/%s.jpg", userId)
	if _, err := os.Stat(imagePath); err == nil {
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.File(imagePath)
		return
	}

	c.File("images/users/user.png")
}

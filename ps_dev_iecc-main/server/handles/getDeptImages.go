package handles

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func GetDeptImage(c *gin.Context) {
	image := c.Param("id")
	imagePath := fmt.Sprintf("%s/%s/%s", "images", "departments", image)

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		c.File("images/default-department.png")
		return
	}

	c.File(imagePath)
}
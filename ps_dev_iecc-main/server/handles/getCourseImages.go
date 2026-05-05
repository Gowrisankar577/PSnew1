package handles

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func GetCourseImage(c *gin.Context) {
	image := c.Param("id")
	imagePath := fmt.Sprintf("%s/%s/%s", "images", "courses", image)

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		c.File("images/default-course.png")
		return
	}

	c.File(imagePath)
}
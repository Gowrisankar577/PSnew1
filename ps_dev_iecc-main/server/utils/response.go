package utils

import "github.com/gin-gonic/gin"

func Response(c *gin.Context, statusCode int, error string, data any) {
	var response map[string]interface{}

	if statusCode != 200 {
		response = map[string]interface{}{
			"success": false,
			"error":   error,
		}
	} else {
		response = map[string]interface{}{
			"success": true,
			"data":    data,
		}
	}

	c.JSON(statusCode, response)
}

package utils

import "github.com/gin-gonic/gin"

func GetClientIP(c *gin.Context) string {
	ip := c.GetHeader("X-Real-Ip")
	if ip == "" {
		ip = c.ClientIP()
	}
	return ip
}

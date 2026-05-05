package handles

import (
	"context"
	"net/http"
	"os"
	"ps_portal/db"
	"ps_portal/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(rate.Limit(1000), 2000)

type AccessLog struct {
	UserID string
	Route  string
	IP     string
}


func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too Many Requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func CorsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowedOrigins := os.Getenv("Origin")
		requestOrigin := c.Request.Header.Get("Origin")

		// Check if request origin is in allowed origins list
		var responseOrigin string
		if allowedOrigins != "" {
			allowedList := strings.Split(allowedOrigins, ",")
			for _, allowed := range allowedList {
				if strings.TrimSpace(allowed) == requestOrigin {
					responseOrigin = requestOrigin
					break
				}
			}
			if responseOrigin == "" {
				// If no match, use the first allowed origin
				responseOrigin = strings.TrimSpace(strings.Split(allowedOrigins, ",")[0])
			}
		} else {
			responseOrigin = requestOrigin
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", responseOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.Writer.WriteHeader(http.StatusOK)
			c.Abort()
			return
		}

		c.Next()
	}
}

func StrictOriginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestOrigin := c.Request.Header.Get("Origin")
		allowedOrigins := os.Getenv("Origin")

		// If no origin header, allow it to proceed
		if requestOrigin == "" {
			c.Next()
			return
		}

		// If no origins configured, deny by default
		if allowedOrigins == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}

		// Check if request origin is in allowed list (supports comma-separated origins)
		allowedList := strings.Split(allowedOrigins, ",")
		for _, allowed := range allowedList {
			if strings.TrimSpace(allowed) == requestOrigin {
				c.Next()
				return
			}
		}

		// Origin not in allowed list
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
	}
}

func ScopeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appBasePath := os.Getenv("APP_BASE_PATH")
		authHeader, _ := c.Cookie("PS")

		user, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized - No user claims found"})
			return
		}

		userDetails, ok := user.(*utils.Claims)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized - Invalid user claims"})
			return
		}

		requestUrl := c.Request.URL.Path
		if appBasePath != "" && strings.HasPrefix(requestUrl, appBasePath) {
			requestUrl = requestUrl[len(appBasePath):]
		}

	
		var count int
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		query := `
			SELECT COUNT(*)
			FROM master_resource_v2 r
			JOIN master_roles ro ON FIND_IN_SET(r.res_group, ro.resources) > 0
			JOIN master_user mu ON mu.role = ro.id
			WHERE r.path = ? AND r.status = '1' 
			AND ro.id = ? AND api_for LIKE '%api%'
			AND mu.id = ? AND mu.status NOT IN ('0', '9') AND mu.auth_token = ?`

		err := db.DB.QueryRowContext(ctx, query, requestUrl, userDetails.Role, userDetails.UserId, authHeader).Scan(&count)
		if err != nil {
			utils.Logging(c, err, 500)
			c.Abort()
			return
		}

		if count == 0 {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized  - Access denied"})
			return
		}

		c.Next()
	}
}
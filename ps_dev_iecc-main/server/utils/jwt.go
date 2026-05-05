package utils

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var JWTKey = []byte("vuovefunepowfefeunfnewupnewnfepnw")

type Claims struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	UserId    string `json:"user_id"`
	UserOffId string `json:"user_off_id"`
	Role      int    `json:"role_id"`
	Dept      string `json:"dept"`
	Year      string `json:"year"`
	YearGroup string `json:"year_group"`
	jwt.StandardClaims
}

// GenerateJWT creates a new JWT token
func GenerateJWT(username, email, id, userOffId, dept, year string, yearGroup string, role int) (string, error) {
	expirationTime := time.Now().Add(3 * time.Hour)
	claims := &Claims{
		Username:  username,
		Email:     email,
		UserId:    id,
		UserOffId: userOffId,
		Role:      role,
		Dept:      dept,
		Year:      year,
		YearGroup: yearGroup,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWTKey)

	return tokenString, err
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		authHeader, err := c.Cookie("PS")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// Parse the token
		token, claims, err := ParseToken(authHeader)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		
		// Set the claims in context for further use
		c.Set("user", claims)
		c.Set("userId", claims.UserId)
		if claims.UserOffId != "" {
			c.Set("UserOffId", claims.UserOffId)
		}else{
			c.Set("UserOffId", claims.UserId)
		}
		c.Set("userName", claims.Username)
		c.Set("userYear", claims.YearGroup)
		c.Set("roleId", claims.Role)

		// Continue with the request
		c.Next()
	}
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JWTKey, nil
	})

	return token, claims, err
}

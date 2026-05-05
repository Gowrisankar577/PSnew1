package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"ps_portal/config"
	"ps_portal/service"
	"ps_portal/utils"

	"cloud.google.com/go/auth/credentials/idtoken"
	"github.com/gin-gonic/gin"
)

type GoogleSignInRequest struct {
	IdToken string `json:"id_token" binding:"required"`
}

func GoogleLogin(c *gin.Context) {
	var req GoogleSignInRequest
	appDomain := os.Getenv("APP_DOMAIN")

	// Bind JSON request to struct and validate required fields
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	payload, err := idtoken.Validate(c, req.IdToken, config.GoogleClientID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Request"})
		return
	}

	// Extract user info from the token
	email, ok := payload.Claims["email"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email not found in token"})
		return
	}

	// Optionally: Create a user in your DB based on Google user info
	// e.g., userService.CreateUser(email, name)

	id, userId, name, dept, year, yearGroup, authToken, role, err := service.LoginService(email, c)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Email ID"})
		} else {
			fmt.Print(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to login, please try again later"})
		}
		return
	}

	token := ""
	if authToken == "" {
		token, err = utils.GenerateJWT(name, email, id, userId, dept, year, yearGroup, role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to login, please try again later"})
			return
		}
	} else {
		token = authToken
		tokenStatus, _, err := utils.ParseToken(authToken)
		if err != nil || !tokenStatus.Valid {
			token, err = utils.GenerateJWT(name, email, id, userId, dept, year, yearGroup, role)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to login, please try again later"})
				return
			}
		}
	}

	service.UpdateAuthToken(id, token)

	c.SetCookie("PS", token, 10800, "/", appDomain, true, true)
	c.JSON(http.StatusOK, gin.H{
		"Message": "Login Success",
		"name": payload.Claims["name"].(string),
	})
}

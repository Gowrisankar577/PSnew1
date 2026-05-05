package service

import (
	"errors"
	"ps_portal/db"

	"github.com/gin-gonic/gin"
)

type UserDetails struct {
	Id                 string `json:"id"`
	UserId             string `json:"user_id"`
	Name               string `json:"name"`
	Dept               string `json:"dept"`
	Year               string `json:"year"`
	YearGroup          string `json:"year_group"`
	Role               int    `json:"role"`
	AllowMultipleLogin string `json:"allow_multiple_login"`
	AuthToken          string `json:"auth_token"`
	Status             string `json:"status"`
}

func LoginService(email string, c *gin.Context) (string, string, string, string, string, string, string, int, error) {
	var user UserDetails
	user, err := fetchUserDetails(email)
	if err != nil {
		return "", "", "", "", "", "", "", 0, err
	}

	if user.Status == "2" {
		logTestSession(c, user.Id)
		return "", "", "", "", "", "", "", 0, errors.New("you already have an active assessment session in another window or device")
	}

	return user.Id, user.UserId, user.Name, user.Dept, user.Year, user.YearGroup, user.AuthToken, user.Role, nil
}

func getClientIP(c *gin.Context) string {
	ip := c.GetHeader("X-Real-Ip")
	if ip == "" {
		ip = c.GetHeader("X-Forwarded-For")
	}
	if ip == "" {
		ip = c.ClientIP()
	}
	return ip
}

func fetchUserDetails(email string) (UserDetails, error) {
	var user UserDetails
	err := db.DB.QueryRow("SELECT mu.id,user_id, mu.name,dept,year,year_group,mu.role, mu.status, mr.allow_multiple_login, ifnull(mu.auth_token, '') FROM master_user mu inner join master_roles mr ON mu.role = mr.id WHERE mu.email = ? AND mu.status NOT IN ('0', '9')", email).Scan(&user.Id, &user.UserId, &user.Name, &user.Dept, &user.Year, &user.YearGroup, &user.Role, &user.Status, &user.AllowMultipleLogin, &user.AuthToken)

	if err == nil && user.AllowMultipleLogin != "1" {
		user.AuthToken = ""
	}

	return user, err
}

// logTestSession logs the IP address if a user attempts to start a test already in progress
func logTestSession(c *gin.Context, userId string) {
	ip := getClientIP(c)
	_, _ = db.DB.Exec("INSERT INTO user_login_request_log (user, request_ip) VALUES(?, ?)", userId, ip)
}

// updateAuthToken updates the auth token for the given user ID in the database
func UpdateAuthToken(userID string, authToken string) error {
	_, err := db.DB.Exec("UPDATE master_user SET auth_token = ? WHERE id = ?", authToken, userID)
	return err
}

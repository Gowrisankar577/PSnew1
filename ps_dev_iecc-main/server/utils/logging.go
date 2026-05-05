package utils

import (
	"fmt"
	"ps_portal/db"

	"github.com/gin-gonic/gin"
)

func Logging(c *gin.Context, err error, code int) {
	fmt.Print(err.Error())
	user := c.MustGet("user")
	userDetails, _ := user.(*Claims)
	_, _ = db.DB.Exec("insert into ps_logging.errors (user,code,error,url) values(?,?,?,?)", userDetails.UserId, code, err.Error(), c.Request.URL.String())
	c.JSON(code, map[string]interface{}{"success": false, "error": "Something went wrong"})
}

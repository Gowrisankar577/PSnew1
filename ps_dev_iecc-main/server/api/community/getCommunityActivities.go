package community

import (
	"ps_portal/db"
	"ps_portal/utils"

	"github.com/gin-gonic/gin"
)

type CommunityActivity struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Time        string `json:"time"`
}

func GetCommunityActivities(c *gin.Context) {
	user := c.MustGet("user")
	userDetails, _ := user.(*utils.Claims)

	rows, err := db.DB.Query(
		`SELECT ca.id, ca.title, ca.description, ca.icon,
		        DATE_FORMAT(ca.created_at, '%d %b %Y %H:%i')
		 FROM community_activities ca
		 JOIN community c ON c.id = ca.community_id AND c.status = '1'
		 JOIN community_members cm ON cm.community_id = c.id AND cm.user_id = ?
		 ORDER BY ca.created_at DESC
		 LIMIT 20`,
		userDetails.UserId,
	)
	if err != nil {
		utils.Logging(c, err, 500)
		return
	}
	defer rows.Close()

	activities := []CommunityActivity{}
	for rows.Next() {
		var a CommunityActivity
		if err := rows.Scan(&a.ID, &a.Title, &a.Description, &a.Icon, &a.Time); err != nil {
			utils.Logging(c, err, 500)
			return
		}
		activities = append(activities, a)
	}

	c.JSON(200, activities)
}

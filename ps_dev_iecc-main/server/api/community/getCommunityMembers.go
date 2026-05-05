package community

import (
	"ps_portal/db"
	"ps_portal/utils"

	"github.com/gin-gonic/gin"
)

type CommunityMember struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	Points int    `json:"points"`
}

func GetCommunityMembers(c *gin.Context) {
	user := c.MustGet("user")
	userDetails, _ := user.(*utils.Claims)

	rows, err := db.DB.Query(
		`SELECT mu.id, mu.name, mu.designation, COALESCE(cm.points, 0)
		 FROM community_members cm
		 JOIN community c ON c.id = cm.community_id AND c.status = '1'
		 JOIN community_members cm2 ON cm2.community_id = c.id AND cm2.user_id = ?
		 JOIN master_user mu ON mu.id = cm.user_id AND mu.status NOT IN ('0','9')
		 ORDER BY cm.points DESC
		 LIMIT 50`,
		userDetails.UserId,
	)
	if err != nil {
		utils.Logging(c, err, 500)
		return
	}
	defer rows.Close()

	members := []CommunityMember{}
	for rows.Next() {
		var m CommunityMember
		if err := rows.Scan(&m.ID, &m.Name, &m.Role, &m.Points); err != nil {
			utils.Logging(c, err, 500)
			return
		}
		members = append(members, m)
	}

	c.JSON(200, members)
}

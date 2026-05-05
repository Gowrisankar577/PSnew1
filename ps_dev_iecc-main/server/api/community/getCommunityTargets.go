package community

import (
	"ps_portal/db"
	"ps_portal/utils"

	"github.com/gin-gonic/gin"
)

type CommunityMandate struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type CommunityTargets struct {
	WeeklyTarget  int                `json:"weekly_target"`
	WeeklyCurrent int                `json:"weekly_current"`
	Mandates      []CommunityMandate `json:"mandates"`
}

func GetCommunityTargets(c *gin.Context) {
	user := c.MustGet("user")
	userDetails, _ := user.(*utils.Claims)

	var targets CommunityTargets
	err := db.DB.QueryRow(
		`SELECT ct.weekly_target, ct.weekly_current
		 FROM community_targets ct
		 JOIN community c ON c.id = ct.community_id AND c.status = '1'
		 JOIN community_members cm ON cm.community_id = c.id AND cm.user_id = ?
		 LIMIT 1`,
		userDetails.UserId,
	).Scan(&targets.WeeklyTarget, &targets.WeeklyCurrent)
	if err != nil {
		utils.Logging(c, err, 500)
		return
	}

	rows, err := db.DB.Query(
		`SELECT cm_t.title, cm_t.completed
		 FROM community_mandates cm_t
		 JOIN community c ON c.id = cm_t.community_id AND c.status = '1'
		 JOIN community_members cm ON cm.community_id = c.id AND cm.user_id = ?
		 ORDER BY cm_t.sort_order ASC`,
		userDetails.UserId,
	)
	if err != nil {
		utils.Logging(c, err, 500)
		return
	}
	defer rows.Close()

	targets.Mandates = []CommunityMandate{}
	for rows.Next() {
		var m CommunityMandate
		if err := rows.Scan(&m.Title, &m.Completed); err != nil {
			utils.Logging(c, err, 500)
			return
		}
		targets.Mandates = append(targets.Mandates, m)
	}

	c.JSON(200, targets)
}

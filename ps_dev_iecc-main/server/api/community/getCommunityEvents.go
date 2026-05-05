package community

import (
	"ps_portal/db"
	"ps_portal/utils"

	"github.com/gin-gonic/gin"
)

type CommunityEvent struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Type        string `json:"type"`
}

func GetCommunityEvents(c *gin.Context) {
	user := c.MustGet("user")
	userDetails, _ := user.(*utils.Claims)

	rows, err := db.DB.Query(
		`SELECT ce.id, ce.title, ce.description,
		        DATE_FORMAT(ce.event_date, '%d %b %Y'), ce.type
		 FROM community_events ce
		 JOIN community c ON c.id = ce.community_id AND c.status = '1'
		 JOIN community_members cm ON cm.community_id = c.id AND cm.user_id = ?
		 WHERE ce.event_date >= CURDATE()
		 ORDER BY ce.event_date ASC
		 LIMIT 10`,
		userDetails.UserId,
	)
	if err != nil {
		utils.Logging(c, err, 500)
		return
	}
	defer rows.Close()

	events := []CommunityEvent{}
	for rows.Next() {
		var e CommunityEvent
		if err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.Date, &e.Type); err != nil {
			utils.Logging(c, err, 500)
			return
		}
		events = append(events, e)
	}

	c.JSON(200, events)
}

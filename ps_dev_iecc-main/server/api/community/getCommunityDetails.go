package community

import (
	"ps_portal/db"
	"ps_portal/utils"

	"github.com/gin-gonic/gin"
)

type CommunityDetails struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Icon             string  `json:"icon"`
	EstablishedDate  string  `json:"established_date"`
	TotalPoints      int     `json:"total_points"`
	MemberCount      int     `json:"member_count"`
	Rating           float64 `json:"rating"`
	Reliability      int     `json:"reliability"`
	Quality          int     `json:"quality"`
	Frequency        int     `json:"frequency"`
}

func GetCommunityDetails(c *gin.Context) {
	user := c.MustGet("user")
	userDetails, _ := user.(*utils.Claims)

	var details CommunityDetails
	err := db.DB.QueryRow(
		`SELECT c.id, c.name, c.icon, DATE_FORMAT(c.established_date, '%d %b %Y'),
		        c.total_points, c.member_count, c.rating,
		        c.reliability, c.quality, c.frequency
		 FROM community c
		 JOIN community_members cm ON cm.community_id = c.id
		 WHERE cm.user_id = ? AND c.status = '1'
		 LIMIT 1`,
		userDetails.UserId,
	).Scan(
		&details.ID, &details.Name, &details.Icon, &details.EstablishedDate,
		&details.TotalPoints, &details.MemberCount, &details.Rating,
		&details.Reliability, &details.Quality, &details.Frequency,
	)
	if err != nil {
		utils.Logging(c, err, 500)
		return
	}

	c.JSON(200, details)
}

package resource

import (
	"ps_portal/db"
	"ps_portal/utils"

	"github.com/gin-gonic/gin"
)

type Menu struct {
	ID      int    `json:"id"`
	Path    string `json:"path"`
	Icon    string `json:"icon"`
	Menu    bool   `json:"menu"`
	Name    string `json:"name"`
	Element string `json:"element"`
}

func GetMyResources(c *gin.Context) {
	user := c.MustGet("user")
	userDetails, _ := user.(*utils.Claims)
	data := []Menu{}

	row, err := db.DB.Query("SELECT r.id,r.path,r.icon,r.menu,r.name,r.element FROM master_resource_v2 r inner join master_resource_group rg on rg.id = r.res_group and rg.status ='1' JOIN master_roles ro ON FIND_IN_SET(r.res_group, resources) > 0 where r.status ='1' and ro.id = ? and r.api_for like '%app%' order by r.sort_by", userDetails.Role)
	if err != nil {
		utils.Logging(c, err, 500)
		return
	}

	for row.Next() {
		var temp Menu
		row.Scan(&temp.ID, &temp.Path, &temp.Icon, &temp.Menu, &temp.Name, &temp.Element)
		if !temp.Menu {
			temp.Name = ""
		}
		data = append(data, temp)
	}

	c.JSON(200, map[string]interface{}{
		"resources": data,
		"id":        c.MustGet("UserOffId"),
		"user_id":   c.MustGet("userId"),
		"user_name": c.MustGet("userName"),
	})
}

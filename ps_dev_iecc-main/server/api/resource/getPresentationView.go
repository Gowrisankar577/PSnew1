package resource

import (
	"ps_portal/db"
	"ps_portal/utils"

	"github.com/gin-gonic/gin"
)


func GetMyPresentationView(c *gin.Context) {
	data := []Activity{}
	user := c.MustGet("user")
	userDetails, _ := user.(*utils.Claims)

	row, err := db.DB.Query("SELECT r.name,r.icon,r.path,r.group FROM master_resource_v2 r inner join master_resource_group rg on rg.id = r.res_group and rg.status ='1' JOIN master_roles ro ON FIND_IN_SET(r.res_group, resources) > 0 where r.status ='1' and ro.id = ? and r.api_for like '%app%' and r.activity in ('2','3') order by r.sort_by", userDetails.Role)
	if err != nil {
		utils.Logging(c, err, 500)
		return
	}

	currIndex := 0
	data = append(data, Activity{})

	for row.Next() {
		var temp ActivityRoutes
		row.Scan(&temp.Name, &temp.Icon, &temp.Path, &temp.Group)
		if data[currIndex].Group != temp.Group {
			data = append(data, Activity{
				Group: temp.Group,
			})
			currIndex++
		}
		data[currIndex].Routes = append(data[currIndex].Routes, temp)

	}

	c.JSON(200, data[1:])
}

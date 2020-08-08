package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"be_nms/models/modelsMember"
	"errors"

	"github.com/labstack/echo/v4"
)

func CreateTargetGroup(c echo.Context) (interface{}, error) {
	authorization := c.Request().Header.Get("Authorization")
	jwt := string([]rune(authorization)[7:])
	tokens, _ := DecodeJWT(jwt)
	db := database.Open()
	defer db.Close()
	admin := models.Admin{}
	db.Where("user_id = ? AND system_id = ?", tokens["user_id"], c.FormValue("systemid")).Find(&admin)
	if admin.ID == 0 {
		return nil, errors.New("You not admin in this system.")
	}
	system := models.System{}
	db.Where("id = ?", c.FormValue("systemid")).Find(&system)
	if system.ID == "" {
		return nil, errors.New("Not have system.")
	}
	targetGroup := modelsMember.TargetGroup{TargetGroupName: c.FormValue("targetgroupname"), NumberOfMembers: 0, SystemID: system.ID}
	db.Create(&targetGroup)
	if targetGroup.ID == 0 {
		return nil, errors.New("Create fail.")
	}
	return targetGroup, nil
}
func GetAllTargetGroup(c echo.Context) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	targetGroups := []modelsMember.TargetGroup{}
	db.Where("system_id = ?", c.QueryParam("systemid")).Find(&targetGroups)
	return targetGroups, nil
}

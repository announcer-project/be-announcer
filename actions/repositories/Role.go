package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"be_nms/models/modelsMember"
	"errors"

	"github.com/labstack/echo/v4"
)

func CreateRole(c echo.Context) (interface{}, error) {
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
	if system.ID == 0 {
		return nil, errors.New("Not have system.")
	}
	role := models.Role{RoleName: c.FormValue("rolename"), SystemID: system.ID}
	db.Create(&role)
	if role.ID == 0 {
		return nil, errors.New("Create fail.")
	}
	targetgroup := modelsMember.TargetGroup{TargetGroupName: c.FormValue("rolename"), NumberOfMembers: 0, SystemID: system.ID}
	db.Create(&targetgroup)
	return role, nil
}

func GetAllRole(c echo.Context) (interface{}, error) {
	authorization := c.Request().Header.Get("Authorization")
	jwt := string([]rune(authorization)[7:])
	tokens, _ := DecodeJWT(jwt)
	db := database.Open()
	defer db.Close()
	system := models.System{}
	db.Where("id = ? AND system_name = ?", c.QueryParam("systemid"), c.QueryParam("systemname")).Find(&system)
	if system.ID == 0 {
		return nil, errors.New("Not have this system.")
	}
	admin := models.Admin{}
	db.Where("user_id = ? AND system_id = ?", tokens["user_id"], c.QueryParam("systemid")).Find(&admin)
	if admin.ID == 0 {
		return nil, errors.New("You not admin in this system.")
	}
	roleuser := []models.Role{}
	db.Where("system_id = ?", c.QueryParam("systemid")).Find(&roleuser)
	return roleuser, nil
}

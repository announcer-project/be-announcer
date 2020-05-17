package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"be_nms/models/modelsMember"
	"be_nms/models/modelsNews"
	"errors"
	"log"

	"github.com/labstack/echo/v4"
)

type User struct {
	Fname          string
	Lname          string
	RoleID         uint
	NewsInterested []modelsNews.NewsType
	Email          string
	Line           string
	SystemID       string
}

func RegisterGetNews(c echo.Context) (interface{}, error) {
	data := User{}
	if err := c.Bind(&data); err != nil {
		return nil, err
	}
	log.Print("Data: ", data)
	db := database.Open()
	defer db.Close()

	tx := db.Begin()

	user := models.User{}
	db.Where("email = ? OR line_id = ?", data.Email, data.Line).First(&user)
	if user.ID == "" {
		user.CreateUser(data.Fname, data.Lname, data.Email, data.Line, "", "")
		tx.Create(&user)
	}
	if user.ID == "" {
		tx.Rollback()
		return nil, errors.New("Register fail.")
	}
	system := models.System{}
	db.Where("id = ?", data.SystemID).Find(&system)
	role := models.Role{}
	db.Where("id = ?", data.RoleID).Find(&role)
	member := modelsMember.Member{UserID: user.ID, SystemID: system.ID, RoleID: role.ID}
	tx.Create(&member)
	if member.ID == 0 {
		tx.Rollback()
		return nil, errors.New("Fail.")
	}
	targetgroup := modelsMember.TargetGroup{}
	db.Where("target_group_name = ?", role.RoleName).Find(&targetgroup)
	if targetgroup.ID == 0 {
		return nil, errors.New("Add to target group error.")
	}
	membergroup := modelsMember.MemberGroup{MemberID: member.ID, TargetGroupID: targetgroup.ID}
	tx.Create(&membergroup)
	if membergroup.ID == 0 {
		return nil, errors.New("Add member to group error.")
	}
	tx.Commit()
	return membergroup, nil
}

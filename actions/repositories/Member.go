package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"be_nms/models/modelsMember"
	"be_nms/models/modelsNews"
	"errors"

	"github.com/labstack/echo/v4"
)

type User struct {
	Fname          string
	Lname          string
	RoleID         uint
	NewsInterested []modelsNews.NewsType
}

func RegisterGetNews(c echo.Context) (interface{}, error) {
	data := User{}
	if err := c.Bind(&data); err != nil {
		return nil, err
	}
	db := database.Open()
	defer db.Close()
	user := models.User{}
	db.Where("email = ? OR line_id = ?", c.FormValue("email"), c.FormValue("line")).First(&user)
	if user.ID != "" {
		return user, errors.New("You have account.")
	}
	user.CreateUser(data.Fname, data.Lname, c.FormValue("email"), c.FormValue("line"), "", "")
	db.Create(&user)
	if user.ID == "" {
		return nil, errors.New("Register fail.")
	}
	system := models.System{}
	db.Where("id = ?", c.FormValue("systemid")).Find(&system)
	role := models.Role{}
	db.Where("id = ?", data.RoleID).Find(&role)
	member := modelsMember.Member{UserID: user.ID, SystemID: system.ID, RoleID: role.ID}
	db.Create(&member)
	if member.ID != 0 {
		return nil, errors.New("Fail.")
	}
	return member, nil
}

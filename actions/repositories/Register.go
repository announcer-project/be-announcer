package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"errors"

	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) (interface{}, error) {
	user := models.User{}
	db := database.Open()
	db.Where("email = ? OR line_id = ?", c.FormValue("email"), c.FormValue("line")).First(&user)
	if user.ID != "" {
		return user, errors.New("You have account.")
	}
	user.CreateUser(c.FormValue("fname"), c.FormValue("lname"), c.FormValue("email"), c.FormValue("line"), c.FormValue("facebook"), c.FormValue("google"))
	db.Create(&user)
	if user.ID == "" {
		return nil, errors.New("Register fail.")
	}
	return user, nil
}

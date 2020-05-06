package handlers

import (
	"be_nms/database"
	"be_nms/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	user := models.User{}
	db := database.Open()
	db.Where("email", c.FormValue("email")).First(&user)
	if user.ID == "" {
		user.CreateUser(c.FormValue("fname"), c.FormValue("lname"), c.FormValue("email"), c.FormValue("line"), c.FormValue("facebook"), c.FormValue("google"))
		db.Create(&user)
		return c.JSON(http.StatusOK, "Create Success")
	}
	return c.JSON(http.StatusBadRequest, "You have account")
}

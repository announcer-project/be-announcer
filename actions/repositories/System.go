package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"errors"
	"log"

	"github.com/labstack/echo/v4"
)

func GetAllsystems(c echo.Context) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	authorization := c.Request().Header.Get("Authorization")
	log.Print(authorization)
	jwt := string([]rune(authorization)[7:])
	tokens, _ := DecodeJWT(jwt)
	admins := []models.Admin{}
	db.Where("user_id = ?", tokens["user_id"]).Find(&admins)
	systems := []models.System{}
	for _, admin := range admins {
		system := models.System{}
		db.Where("id = ?", admin.SystemID).First(&system)
		systems = append(systems, system)
	}
	return systems, nil
}

func CreateSystem(c echo.Context) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	authorization := c.Request().Header.Get("Authorization")
	jwt := string([]rune(authorization)[7:])
	tokens, _ := DecodeJWT(jwt)
	user := models.User{}
	db.Where("id = ?", tokens["user_id"]).Find(&user)
	if user.ID == "" {
		return nil, errors.New("Create fail.")
	}
	system := models.System{SystemName: c.FormValue("systemname"), OwnerID: user.ID}
	db.Create(&system)
	db.First(&system)
	if system.ID == 0 {
		return nil, errors.New("Create fail.")
	}

	admin := models.Admin{UserID: user.ID, SystemID: system.ID, Position: "admin"}
	db.Create(&admin)
	db.First(&admin)
	if system.ID == 0 {
		return nil, errors.New("Create fail.")
	}
	return user, nil
}

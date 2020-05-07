package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"log"

	"github.com/labstack/echo/v4"
)

func GetAllsystems(c echo.Context) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	authorization := c.Request().Header.Get("Authorization")
	jwt := string([]rune(authorization)[7:])
	tokens, _ := DecodeJWT(jwt)
	admins := []models.Admin{}
	log.Print("userid",tokens["user_id"])
	db.Where("user_id = ?",tokens["user_id"]).Find(&admins)
	systems := []models.System{}
	system := models.System{}
	for i, admin := range admins {
		db.Where("id = ?", admin.SystemID).First(&system)
		log.Print(system)
		log.Print(i)
		systems = append(systems, system)
	}
	return systems, nil
}

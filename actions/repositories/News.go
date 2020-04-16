package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"time"

	"github.com/labstack/echo/v4"
)

func CreateNews(c echo.Context) (bool, error) {
	title := c.FormValue("title")
	content := c.FormValue("content")
	expireDate, _ := time.Parse(c.FormValue("expiredate"), "YYYY-MM-DD")
	jwt := c.FormValue("jwt")
	tokens, _ := DecodeJWT(jwt)
	db := database.Open()
	defer db.Close()
	admin := models.Admin{}
	db.Where("user_id = ?", tokens["user_id"]).First(&admin)
	news := models.News{}
	news.CreateNews(title, content, expireDate, admin.AdminID)
	db.Create(&news)
	return true, nil
}

func GetNewsByID(c echo.Context) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	news := models.News{}
	if c.Param("id") != "" {
		db.First(&news, c.Param("id"))
	} else {
		db.First(&news, c.FormValue("id"))
	}
	return news, nil
}

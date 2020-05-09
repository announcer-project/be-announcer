package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"be_nms/models/modelsNews"
	"errors"
	"log"
	"time"

	"github.com/labstack/echo/v4"
)

func CreateNews(c echo.Context) (bool, error) {
	title := c.FormValue("title")
	body := c.FormValue("body")
	expireDate, _ := time.Parse(c.FormValue("expiredate"), "YYYY-MM-DD")
	authorization := c.Request().Header.Get("Authorization")
	jwt := string([]rune(authorization)[7:])
	tokens, _ := DecodeJWT(jwt)
	db := database.Open()
	defer db.Close()
	admin := models.Admin{}
	db.Where("user_id = ? AND system_id = ?", tokens["user_id"], c.FormValue("systemid")).Find(&admin)
	log.Print(admin)
	log.Print(tokens["user_id"])
	log.Print(c.FormValue("systemid"))
	system := models.System{}
	db.Where("id = ?", c.FormValue("systemid")).First(&system)
	news := modelsNews.News{Title: title, Body: body, ExpireDate: expireDate, SystemID: system.ID, AuthorID: admin.ID}
	db.Create(&news)
	return true, nil
}

func GetNewsByID(c echo.Context) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	news := modelsNews.News{}
	if c.Param("id") != "" {
		db.First(&news, c.Param("id"))
	} else {
		db.First(&news, c.FormValue("id"))
	}
	return news, nil
}

func GetAllNews(c echo.Context) (interface{}, error) {
	authorization := c.Request().Header.Get("Authorization")
	jwt := string([]rune(authorization)[7:])
	tokens, _ := DecodeJWT(jwt)
	db := database.Open()
	defer db.Close()
	admin := models.Admin{}
	db.Where("user_id = ? AND system_id = ?", tokens["user_id"], c.FormValue("systemid")).Find(&admin)
	log.Print(admin)
	log.Print(tokens["user_id"])
	log.Print(c.FormValue("test"))
	log.Print(c.FormValue("systemid"))
	if admin.ID == 0 {
		return nil, errors.New("You not admin in this system.")
	}
	news := []modelsNews.News{}
	db.Where("system_id = ?", c.FormValue("systemid")).Find(&news)
	return news, nil
}

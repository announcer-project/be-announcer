package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"be_nms/models/modelsNews"
	"errors"
	"time"

	"github.com/labstack/echo/v4"
)

type Data struct {
	Title      string
	Body       string
	Expiredate string
	NewsTypes  []string
	SystemID   uint
}

//News
func CreateNews(c echo.Context) (bool, error) {
	data := Data{}
	if err := c.Bind(&data); err != nil {
		return false, err
	}
	authorization := c.Request().Header.Get("Authorization")
	jwt := string([]rune(authorization)[7:])
	tokens, _ := DecodeJWT(jwt)
	db := database.Open()
	defer db.Close()
	admin := models.Admin{}
	db.Where("user_id = ? AND system_id = ?", tokens["user_id"], data.SystemID).Find(&admin)
	system := models.System{}
	db.Where("id = ?", data.SystemID).First(&system)
	if system.ID == 0 {
		return false, errors.New("Create fail.")
	}
	expiredate, _ := time.Parse("dd-mm-yy", data.Expiredate)
	news := modelsNews.News{Title: data.Title, Body: data.Body, ExpireDate: expiredate, SystemID: system.ID, AuthorID: admin.ID}
	db.Create(&news)
	if news.ID == 0 {
		return false, errors.New("Create fail.")
	}
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
	if admin.ID == 0 {
		return nil, errors.New("You not admin in this system.")
	}
	news := []modelsNews.News{}
	db.Where("system_id = ?", c.FormValue("systemid")).Find(&news)
	return news, nil
}

//NewsType
func CreateNewsType(c echo.Context) (interface{}, error) {
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
	newsType := modelsNews.NewsType{NewTypeName: c.FormValue("newstypename"), SystemID: system.ID}
	db.Create(&newsType)
	if newsType.ID == 0 {
		return nil, errors.New("Create fail.")
	}
	return newsType, nil
}
func GetAllNewsType(c echo.Context) (interface{}, error) {
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
	newsTypes := []modelsNews.NewsType{}
	db.Where("system_id = ?", c.FormValue("systemid")).Find(&newsTypes)
	return newsTypes, nil
}

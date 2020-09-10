package handlers

import (
	"be_nms/actions/repositories"
	"be_nms/models"
	"be_nms/models/modelsNews"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

//News
func CreateNews(c echo.Context) error {
	newsID, err := repositories.CreateNews(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, newsID)
}

func GetNewsByID(c echo.Context) error {
	news, _ := repositories.GetNewsByID(c, c.Param("id"))
	return c.JSON(http.StatusOK, news)
}

// func GetAllNews(c echo.Context) error {
// 	news, _ := repositories.GetAllNews(c)
// 	return c.JSON(http.StatusOK, news)
// }
type AllNewsClassify struct {
	NewsDraft   []modelsNews.News `json:"newsdraft"`
	NewsPublish []modelsNews.News `json:"newspublish"`
}

func GetAllNewsByClassify(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	newsDraft, _ := repositories.GetAllNews(tokens["user_id"].(string), c.QueryParam("systemid"), "draft")
	newsPublish, _ := repositories.GetAllNews(tokens["user_id"].(string), c.QueryParam("systemid"), "publish")
	allnews := AllNewsClassify{NewsDraft: newsDraft.([]modelsNews.News), NewsPublish: newsPublish.([]modelsNews.News)}
	return c.JSON(http.StatusOK, allnews)
}
func GetAllNewsDraft(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	news, _ := repositories.GetAllNews(tokens["user_id"].(string), c.QueryParam("systemid"), "draft")
	return c.JSON(http.StatusOK, news)
}
func GetAllNewsPublish(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	news, _ := repositories.GetAllNews(tokens["user_id"].(string), c.QueryParam("systemid"), "publish")
	return c.JSON(http.StatusOK, news)
}

//NewsType
func CreateNewsType(c echo.Context) error {
	_, err := repositories.CreateNewsType(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, "Create Success.")
}

func GetAlNewsType(c echo.Context) error {
	newsTypes, err := repositories.GetAllNewsType(c.QueryParam("systemid"), true)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, newsTypes)
}

func DeleteNewsType(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	var data struct {
		Systemid   string
		Newstypeid int
	}
	if err := c.Bind(&data); err != nil {
		log.Print("error ", err)
		return err
	}
	log.Print(tokens["user_id"])
	log.Print(data)
	if err := repositories.DeleteNewsType(tokens["user_id"].(string), data.Systemid, data.Newstypeid); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, "Test")
}

//Announce
func AnnounceNews(c echo.Context) error {
	log.Print(c.FormValue("newsid"))
	log.Print(c.FormValue("systemid"))
	news, err := repositories.GetNewsByID(c, c.FormValue("newsid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	log.Print(news)
	system, err := repositories.GetSystemByID(c, c.FormValue("systemid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	announce, err := repositories.BroadcastNewsLine(c, news.(modelsNews.News), system.(models.System))
	if !announce {
		return c.JSON(http.StatusOK, err)
	}
	return c.JSON(http.StatusOK, "Announce Success!")
}

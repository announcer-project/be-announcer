package handlers

import (
	"be_nms/actions/repositories"
	"be_nms/database"
	"be_nms/models"
	"be_nms/models/modelsNews"
	"net/http"

	"github.com/labstack/echo/v4"
)

//News
func CreateNews(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	var data struct {
		Cover           string
		Title           string
		Body            string
		CheckExpiredate bool
		Expiredate      string
		Images          []string
		NewsTypes       []struct {
			ID   int
			Name string
		}
		SystemID string
		Status   string
	}
	if err := c.Bind(&data); err != nil {
		message.Message = "server error."
		return c.JSON(500, message)
	}
	newsID, err := repositories.CreateNews(
		data.Cover,
		data.SystemID,
		tokens["user_id"].(string),
		data.CheckExpiredate,
		data.Expiredate,
		data.Title,
		data.Body,
		data.Status,
		data.NewsTypes,
		data.Images)
	if err != nil {
		message.Message = err.Error()
		return c.JSON(http.StatusBadRequest, message)
	}
	var success struct {
		NewsID uint
	}
	success.NewsID = newsID.(uint)
	return c.JSON(http.StatusOK, success)
}

func GetNewsByID(c echo.Context) error {
	var message struct {
		Message string `json:"message"`
	}
	news, err := repositories.GetNewsByID("publish", c.Param("id"))
	if err != nil {
		message.Message = err.Error()
		return c.JSON(400, message)
	}
	return c.JSON(http.StatusOK, news)
}

func GetNewsDraftByID(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	if c.QueryParam("systemid") == "" {
		message.Message = "not have query param"
		return c.JSON(400, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	admin := models.Admin{}
	db := database.Open()
	defer db.Close()
	db.Where("user_id = ? AND system_id = ?", tokens["user_id"].(string), c.QueryParam("systemid")).Find(&admin)
	if admin.ID == 0 {
		message.Message = "You not admin in this system."
		return c.JSON(400, message)
	}
	news, err := repositories.GetNewsByID("draft", c.Param("id"))
	if err != nil {
		message.Message = err.Error()
		return c.JSON(400, message)
	}
	return c.JSON(http.StatusOK, news)
}

type AllNewsClassify struct {
	NewsDraft   []modelsNews.News `json:"newsdraft"`
	NewsPublish []modelsNews.News `json:"newspublish"`
}

func GetAllNewsByClassify(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	if c.QueryParam("systemid") == "" {
		message.Message = "not have query param systemid."
		return c.JSON(400, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	newsDraft, err := repositories.GetAllNews(tokens["user_id"].(string), c.QueryParam("systemid"), "draft")
	if err != nil {
		message.Message = err.Error()
		return c.JSON(401, err)
	}
	newsPublish, _ := repositories.GetAllNews(tokens["user_id"].(string), c.QueryParam("systemid"), "publish")
	if err != nil {
		message.Message = err.Error()
		return c.JSON(401, err)
	}
	allnews := struct {
		NewsDraft   []modelsNews.News `json:"draft"`
		NewsPublish []modelsNews.News `json:"publish"`
	}{
		newsDraft.([]modelsNews.News),
		newsPublish.([]modelsNews.News),
	}
	return c.JSON(http.StatusOK, allnews)
}
func GetAllNewsDraft(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	if c.QueryParam("systemid") == "" {
		message.Message = "not have query param systemid."
		return c.JSON(400, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	news, err := repositories.GetAllNews(tokens["user_id"].(string), c.QueryParam("systemid"), "draft")
	if err != nil {
		message.Message = err.Error()
		return c.JSON(401, err)
	}
	return c.JSON(http.StatusOK, news)
}
func GetAllNewsPublish(c echo.Context) error {
	var message struct {
		Message string `json:"message"`
	}
	if c.QueryParam("systemid") == "" {
		message.Message = "not have query param systemid."
		return c.JSON(400, message)
	}
	news, err := repositories.GetAllNews("publish", c.QueryParam("systemid"), "publish")
	if err != nil {
		message.Message = err.Error()
		return c.JSON(401, err)
	}
	return c.JSON(http.StatusOK, news)
}

//NewsType
func CreateNewsType(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	var data struct {
		SystemID     string
		NewsTypeName string
	}
	if err := c.Bind(&data); err != nil {
		message.Message = "server error."
		return c.JSON(500, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	_, err := repositories.CreateNewsType(tokens["user_id"].(string), data.SystemID, data.NewsTypeName)
	if err != nil {
		message.Message = err.Error()
		return c.JSON(400, message)
	}
	message.Message = "create success."
	return c.JSON(http.StatusOK, message)
}

func GetAlNewsType(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	if c.QueryParam("systemid") == "" {
		message.Message = "not have query param."
		return c.JSON(401, message)
	}
	newsTypes, err := repositories.GetAllNewsType(c.QueryParam("systemid"), true)
	if err != nil {
		message.Message = err.Error()
		return c.JSON(401, message)
	}
	return c.JSON(http.StatusOK, newsTypes)
}

func DeleteNewsType(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	var data struct {
		Systemid   string
		NewsTypeid int
	}
	if err := c.Bind(&data); err != nil {
		message.Message = "server error."
		return c.JSON(500, message)
	}
	if err := repositories.DeleteNewsType(tokens["user_id"].(string), data.Systemid, data.NewsTypeid); err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	message.Message = "delete success."
	return c.JSON(http.StatusOK, message)
}

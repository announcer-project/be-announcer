package openapi

import (
	"be_nms/actions/repositories"
	"be_nms/actions/repositories/openapi"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetNewsbyID(c echo.Context) error {
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
	if c.Param("id") == "" {
		message.Message = "not have param"
		return c.JSON(400, message)
	}
	news, err := openapi.GetNewsByID(c.Param("id"), c.QueryParam("systemid"))
	if err != nil {
		message.Message = err.Error()
		return c.JSON(400, message)
	}
	return c.JSON(200, news)
}

func GetAllNewsPublish(c echo.Context) error {
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
	news := openapi.GetAllNews("publish", c.QueryParam("systemid"))
	return c.JSON(200, news)
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
		message.Message = "not have query param"
		return c.JSON(400, message)
	}
	news := openapi.GetAllNews("draft", c.QueryParam("systemid"))
	return c.JSON(200, news)
}

func CreateNews(c echo.Context) error {
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
		Status string
	}
	if err := c.Bind(&data); err != nil {
		message.Message = "server error."
		return c.JSON(500, message)
	}
	newsID, err := repositories.CreateNews(
		data.Cover,
		c.QueryParam("systemid"),
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

func DeleteNews(c echo.Context) error {
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

	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	var data struct {
		Newsid int
	}
	if err := c.Bind(&data); err != nil {
		message.Message = "server error."
		return c.JSON(500, message)
	}
	if err := repositories.DeleteNews(tokens["user_id"].(string), c.QueryParam("systemid"), data.Newsid); err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	message.Message = "delete success."
	return c.JSON(http.StatusOK, message)
}

func CreateNewsType(c echo.Context) error {
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

	var data struct {
		NewsTypeName string
	}
	if err := c.Bind(&data); err != nil {
		message.Message = "server error."
		return c.JSON(500, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	_, err := repositories.CreateNewsType(tokens["user_id"].(string), c.QueryParam("systemid"), data.NewsTypeName)
	if err != nil {
		message.Message = err.Error()
		return c.JSON(400, message)
	}
	message.Message = "create success."
	return c.JSON(http.StatusOK, message)
}

func GetAllNewsType(c echo.Context) error {
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

	if c.QueryParam("systemid") == "" {
		message.Message = "not have query param."
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
	if err := repositories.DeleteNewsType(tokens["user_id"].(string), c.QueryParam("systemid"), data.NewsTypeid); err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	message.Message = "delete success."
	return c.JSON(http.StatusOK, message)
}

package handlers

import (
	"be_nms/actions/repositories"
	"be_nms/models/modelsNews"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateMember(c echo.Context) error {
	var data struct {
		IsUser         bool
		FName          string
		LName          string
		Email          string
		ImageUrl       string
		RoleID         int
		NewsInterested []modelsNews.NewsType
		SystemID       string
		Line           string
	}
	var message struct {
		Message string `json:"message"`
	}
	if err := c.Bind(&data); err != nil {
		message.Message = "server error."
		return c.JSON(500, message)
	}
	log.Print(data)
	err := repositories.RegisterGetNews(data)
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	message.Message = "register for get news success."
	return c.JSON(http.StatusOK, message)
}

func GetAllMember(c echo.Context) error {
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
		return c.JSON(400, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	members, err := repositories.GetAllMember(tokens["user_id"].(string), c.QueryParam("systemid"))
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	return c.JSON(http.StatusOK, members)
}

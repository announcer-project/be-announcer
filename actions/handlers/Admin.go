package handlers

import (
	"be_nms/actions/repositories"

	"github.com/labstack/echo/v4"
)

func GetAllAdmin(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	if c.Param("systemid") == "" {
		message.Message = "not have param."
		return c.JSON(400, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	admins, err := repositories.GetAllAdmin(c.Param("systemid"), tokens["user_id"].(string))
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	return c.JSON(200, admins)
}

func CreateAdmin(c echo.Context) error {
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
	var data struct {
		UserID string
	}
	if err := c.Bind(&data); err != nil {
		message.Message = "server error."
		return c.JSON(500, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	err := repositories.CreateAdmin(c.QueryParam("systemid"), tokens["user_id"].(string), data.UserID)
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	message.Message = "add co-admin success."
	return c.JSON(200, message)
}

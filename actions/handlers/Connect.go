package handlers

import (
	"be_nms/actions/repositories"

	"github.com/labstack/echo/v4"
)

func CheckConnectLineOA(c echo.Context) error {
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
	check := repositories.CheclConnectLineOA(c.QueryParam("systemid"))
	return c.JSON(200, check)
}

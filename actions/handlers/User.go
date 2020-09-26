package handlers

import (
	"be_nms/actions/repositories"

	"github.com/labstack/echo/v4"
)

func GetUserByJWT(c echo.Context) error {
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
	user, err := repositories.GetUserByID(tokens["user_id"].(string))
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	return c.JSON(200, user)
}

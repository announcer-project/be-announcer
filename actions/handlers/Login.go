package handlers

import (
	"be_nms/actions/repositories"
	"be_nms/models"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	social := c.Request().Header.Get("Social")
	userID := ""
	err := errors.New("error")
	if social == "line" {
		userID, err = repositories.GetUserIDLine(c)
		if err != nil {
			return c.JSON(400, err)
		}
	} else if social == "facebook" {
		userID = c.Request().Header.Get("UserID")
	}
	user, err := repositories.GetUserBySocialId(userID, social)
	if err != nil {
		return c.JSON(http.StatusBadRequest, userID)
	}
	jwt := repositories.EncodeJWT(user.(models.User))
	return c.JSON(http.StatusOK, jwt)
}

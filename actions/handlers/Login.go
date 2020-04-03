package handlers

import (
	"be_nms/actions/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

func LineLogin(c echo.Context) error {
	token, err := repositories.GetAccessTokenLine(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}
	userId, err := repositories.GetUserProfileLine(token)
	return c.JSON(http.StatusOK, userId)
}

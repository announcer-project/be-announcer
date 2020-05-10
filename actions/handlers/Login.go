package handlers

import (
	"be_nms/actions/repositories"
	"be_nms/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func LineLogin(c echo.Context) error {
	userIDLine, err := repositories.GetUserIDLine(c)
	if err != nil {
		return c.JSON(400, err)
	}
	user, err := repositories.GetUserBySocialId(userIDLine, "line")
	if err != nil {
		return c.JSON(http.StatusBadRequest, userIDLine)
	}
	jwt := repositories.EncodeJWT(user.(models.User))
	return c.JSON(http.StatusOK, jwt)
}

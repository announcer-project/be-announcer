package handlers

import (
	"be_nms/actions/repositories"
	"be_nms/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	user, err := repositories.Register(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	jwt := repositories.EncodeJWT(user.(models.User))
	return c.JSON(http.StatusOK, jwt)
}

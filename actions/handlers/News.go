package handlers

import (
	"be_nms/actions/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateNews(c echo.Context) error {
	news, _ := repositories.CreateNews(c)
	return c.JSON(http.StatusOK, news)
}

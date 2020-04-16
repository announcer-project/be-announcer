package handlers

import (
	"be_nms/actions/repositories"
	"be_nms/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateNews(c echo.Context) error {
	repositories.CreateNews(c)
	return c.JSON(http.StatusOK, "Create success.")
}

func GetNewsByID(c echo.Context) error {
	news, _ := repositories.GetNewsByID(c)
	return c.JSON(http.StatusOK, news)
}

func AnnounceNews(c echo.Context) error {
	news := models.News{}
	announce := repositories.BroadMessageLine(c, news)
	if !announce {
		return c.JSON(http.StatusOK, "Announe Fail!")
	}
	return c.JSON(http.StatusOK, "Announe Success!")
}

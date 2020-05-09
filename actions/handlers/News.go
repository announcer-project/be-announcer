package handlers

import (
	"be_nms/actions/repositories"
	"be_nms/models/modelsNews"
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

func GetAllNews(c echo.Context) error {
	news, _ := repositories.GetAllNews(c)
	return c.JSON(http.StatusOK, news)
}
func AnnounceNews(c echo.Context) error {
	news, err := repositories.GetNewsByID(c)
	if err != nil {
		return c.JSON(http.StatusOK, err)
	}
	announce, err := repositories.BroadcastNewsLine(c, news.(modelsNews.News))
	if !announce {
		return c.JSON(http.StatusOK, err)
	}
	return c.JSON(http.StatusOK, "Announe Success!")
}

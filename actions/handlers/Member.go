package handlers

import (
	"be_nms/actions/repositories"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateMember(c echo.Context) error {
	_, err := repositories.RegisterGetNews(c)
	if err != nil {
		log.Print("Error : ", err)
		return c.JSON(401, err)
	}
	return c.String(http.StatusOK, "OK!")
}

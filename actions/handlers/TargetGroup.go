package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateTargetGroup(c echo.Context) error {
	return c.JSON(http.StatusOK, "test")
}

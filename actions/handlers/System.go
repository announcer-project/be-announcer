package handlers

import (
	"be_nms/actions/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetAllSystems(c echo.Context) error {
	systems, err := repositories.GetAllsystems(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, systems)
}

func CreateSystem(c echo.Context) error {
	_, err := repositories.CreateSystem(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, "Create Success.")
}
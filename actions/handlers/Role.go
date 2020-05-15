package handlers

import (
	"be_nms/actions/repositories"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateRole(c echo.Context) error {
	_, err := repositories.CreateRole(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, "Create success.")
}

func GetAllRole(c echo.Context) error {
	roleuser, err := repositories.GetAllRole(c)
	if err != nil {
		return err
	}
	log.Print(roleuser)
	return c.JSON(http.StatusOK, roleuser)
}

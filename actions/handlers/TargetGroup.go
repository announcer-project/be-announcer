package handlers

import (
	"be_nms/actions/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateTargetGroup(c echo.Context) error {
	_, err := repositories.CreateTargetGroup(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, "Create success.")
}
func GetAllTargetGroup(c echo.Context) error {
	targetGroups, err := repositories.GetAllTargetGroup(c.QueryParam("systemid"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, targetGroups)
}

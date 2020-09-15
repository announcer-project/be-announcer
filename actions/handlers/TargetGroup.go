package handlers

import (
	"be_nms/actions/repositories"
	"log"
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
	var message struct {
		Message string `json:"message"`
	}
	log.Print(c.QueryParam("systemid"))
	if c.QueryParam("systemid") == "" {
		message.Message = "not have query param."
		return c.JSON(400, message)
	}
	targetGroups, err := repositories.GetAllTargetGroup(c.QueryParam("systemid"))
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	return c.JSON(http.StatusOK, targetGroups)
}

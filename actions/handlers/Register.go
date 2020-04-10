package handlers

import (
	"be_nms/actions/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	status, err := repositories.Register(c)
	if status == false {
		return c.JSON(400, err)
	}
	return c.JSON(http.StatusOK, "Send to email success.")
}

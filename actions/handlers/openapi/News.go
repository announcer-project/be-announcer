package openapi

import (
	"be_nms/actions/repositories/openapi"

	"github.com/labstack/echo/v4"
)

func GetNewsbyID(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	if c.QueryParam("systemid") == "" {
		message.Message = "not have query param"
		return c.JSON(400, message)
	}
	if c.Param("id") == "" {
		message.Message = "not have param"
		return c.JSON(400, message)
	}
	news, err := openapi.GetNewsByID(c.Param("id"), c.QueryParam("systemid"))
	if err != nil {
		message.Message = err.Error()
		return c.JSON(400, message)
	}
	return c.JSON(200, news)
}

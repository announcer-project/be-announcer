package handlers

import (
	"be_nms/actions/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateRole(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	var data struct {
		SystemID string
		Rolename string
		Require  bool
	}
	if err := c.Bind(&data); err != nil {
		message.Message = "server error."
		return c.JSON(500, message)
	}
	_, err := repositories.CreateRole(tokens["user_id"].(string), data.SystemID, data.Rolename, data.Require)
	if err != nil {
		message.Message = err.Error()
		return c.JSON(400, message)
	}
	message.Message = "create role success."
	return c.JSON(http.StatusOK, message)
}

func GetAllRole(c echo.Context) error {
	var message struct {
		Message string `json:"message"`
	}
	if c.QueryParam("systemid") == "" {
		message.Message = "not have query param."
		return c.JSON(400, message)
	}
	roleuser, err := repositories.GetAllRole(c.QueryParam("systemid"))
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	return c.JSON(http.StatusOK, roleuser)
}

func GetRoleRequest(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	if c.Param("systemid") == "" {
		message.Message = "not have param."
		return c.JSON(400, message)
	}
	members, err := repositories.GetRoleRequest(c.Param("systemid"))
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	return c.JSON(200, members)
}

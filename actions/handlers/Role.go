package handlers

import (
	"be_nms/actions/repositories"
	"log"
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

func DeleteRole(c echo.Context) error {
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
	if c.Param("systemid") == "" || c.Param("roleid") == "" {
		message.Message = "not found param."
		return c.JSON(400, message)
	}
	err := repositories.DeleteRole(c.Param("systemid"), tokens["user_id"].(string), c.Param("roleid"))
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	message.Message = "delete success."
	return c.JSON(200, message)
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

func ApproveRoleRequest(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	var data struct {
		MemberID string
		SystemID string
	}
	if err := c.Bind(&data); err != nil {
		message.Message = "server error."
		return c.JSON(500, message)
	}
	log.Print(data)

	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	err := repositories.ApproveRoleRequest(data.MemberID, tokens["user_id"].(string), data.SystemID)
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	message.Message = "approve success."
	return c.JSON(200, message)
}

func RejectRoleRequest(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	log.Print(authorization)
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	var data struct {
		MemberID string
		SystemID string
	}
	if err := c.Bind(&data); err != nil {
		message.Message = "server error."
		return c.JSON(500, message)
	}
	log.Print(data)
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	err := repositories.RejectRoleRequest(data.MemberID, tokens["user_id"].(string), data.SystemID)
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	message.Message = "reject success."
	return c.JSON(200, message)
}

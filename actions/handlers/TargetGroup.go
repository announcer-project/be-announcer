package handlers

import (
	"be_nms/actions/repositories"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateTargetGroup(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "Not have jwt."
		return c.JSON(401, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	var data struct {
		Groupname string
		SystemID  string
		Members   []struct {
			MemberID uint
		}
	}
	if err := c.Bind(&data); err != nil {
		message.Message = "server error."
		return c.JSON(500, message)
	}
	_, err := repositories.CreateTargetGroup(tokens["user_id"].(string), data.SystemID, data.Groupname, data.Members)
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	message.Message = "create target group success."
	return c.JSON(http.StatusOK, message)
}

func GetAllTargetGroup(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "Not have jwt."
		return c.JSON(401, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	if c.Param("systemid") == "" {
		message.Message = "not have query param."
		return c.JSON(400, message)
	}
	targetGroups, err := repositories.GetAllTargetGroup(tokens["user_id"].(string), c.Param("systemid"))
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	return c.JSON(http.StatusOK, targetGroups)
}

func DeleteTargetGroup(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "Not have jwt."
		return c.JSON(401, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	if c.Param("systemid") == "" || c.Param("targetgroupid") == "" {
		message.Message = "not have query param."
		return c.JSON(400, message)
	}
	err := repositories.DeleteTargetgroup(c.Param("systemid"), tokens["user_id"].(string), c.Param("targetgroupid"))
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	message.Message = "delete success"
	return c.JSON(200, message)
}

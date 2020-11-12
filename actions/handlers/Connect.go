package handlers

import (
	"be_nms/actions/repositories"
	"log"

	"github.com/labstack/echo/v4"
)

func CheckConnectLineOA(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	if c.QueryParam("systemid") == "" {
		message.Message = "not have query param."
		return c.JSON(400, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	check, err := repositories.CheckConnectLineOA(c.QueryParam("systemid"), tokens["user_id"].(string))
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	return c.JSON(200, check)
}

func DisconnectLinaOA(c echo.Context) error {
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
		return c.JSON(401, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	err := repositories.DisconnectLineOA(c.Param("systemid"), tokens["user_id"].(string))
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	message.Message = "disconnect success."
	return c.JSON(200, message)
}

func ConenctLineOA(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	var data struct {
		SystemID           string
		ChannelID          string
		ChannelAccessToken string
		LiffID             string
		Roles              []struct {
			Rolename string
			Require  bool
		}
	}
	if err := c.Bind(&data); err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	log.Print(data)
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	connecterr := repositories.ConnectLineOA(data.SystemID, tokens["user_id"].(string), data.ChannelID, data.LiffID, data.ChannelAccessToken, data.Roles)
	if connecterr != nil {
		message.Message = connecterr.Error()
		return c.JSON(500, message)
	}
	message.Message = "connect success."
	return c.JSON(200, message)
}

func GetLiffID(c echo.Context) error {
	var message struct {
		Message string `json:"message"`
	}
	if c.QueryParam("systemid") == "" {
		message.Message = "not have query param."
		return c.JSON(401, message)
	}
	liffid, err := repositories.GetLiffID(c.QueryParam("systemid"))
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	return c.JSON(200, liffid)
}

package handlers

import (
	"be_nms/actions/repositories"

	"github.com/labstack/echo/v4"
)

func CheckConnectDialogflow(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	if c.QueryParam("systemid") == "" {
		message.Message = "Not have system ID."
		return c.JSON(500, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	checked, err := repositories.CheckConnectDialogflow(tokens["user_id"].(string), c.QueryParam("systemid"))
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	return c.JSON(200, checked)
}

func ConnectDialogflow(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	if c.QueryParam("systemid") == "" {
		message.Message = "Not have system ID."
		return c.JSON(500, message)
	}
	var data struct {
		ProjectID          string
		AuthJSONFileBase64 string
		Lang               string
		TimeZone           string
	}
	if err := c.Bind(&data); err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	err := repositories.ConnectDialogflow(tokens["user_id"].(string), c.QueryParam("systemid"), data.ProjectID, data.AuthJSONFileBase64, data.Lang, data.TimeZone)
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	message.Message = "success"
	return c.JSON(200, message)
}

type Event struct {
	ReplyToken string
	Type       string
	Source     struct {
		Type   string
		UserID string
	}
	Message struct {
		ID   string
		Type string
		Text string
	}
}

type MessageEvent struct {
	Destination string
	Events      []Event
}

func Webhook(c echo.Context) error {
	var message struct {
		Message string `json:"message"`
	}
	if c.Param("systemid") == "" {
		message.Message = "Not have system ID."
		return c.JSON(500, message)
	}
	messageEvent := MessageEvent{}
	if err := c.Bind(&messageEvent); err != nil {
		message.Message = "server error."
		return c.JSON(500, message)
	}
	response, err := repositories.Webhook(c.Param("systemid"), messageEvent.Events[0].Message.Text)
	if err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	return c.JSON(200, response)
}

func ListIntent(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	if c.QueryParam("systemid") == "" {
		message.Message = "Not have system ID."
		return c.JSON(500, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	response, err := repositories.ListIntents(tokens["user_id"].(string), c.QueryParam("systemid"))
	if err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, response)
}

func GetIntent(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	if c.QueryParam("systemid") == "" {
		message.Message = "Not have system ID."
		return c.JSON(500, message)
	}
	var data struct {
		IntentName string
	}
	if err := c.Bind(&data); err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	response, err := repositories.GetIntent(tokens["user_id"].(string), c.QueryParam("systemid"), data.IntentName)
	if err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, response)
}

func CreateIntent(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	if c.QueryParam("systemid") == "" {
		message.Message = "Not have system ID."
		return c.JSON(500, message)
	}
	var data struct {
		DisplayName    string
		TrainingPhrase []string
		MessageTexts   []string
	}
	if err := c.Bind(&data); err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	err := repositories.CreateIntent(tokens["user_id"].(string), c.QueryParam("systemid"), data.DisplayName, data.TrainingPhrase, data.MessageTexts)
	if err != nil {
		return c.JSON(500, err)
	}
	message.Message = "success"
	return c.JSON(200, message)
}

func DeleteIntent(c echo.Context) error {
	authorization := c.Request().Header.Get("Authorization")
	var message struct {
		Message string `json:"message"`
	}
	if authorization == "" {
		message.Message = "not have jwt."
		return c.JSON(401, message)
	}
	if c.QueryParam("systemid") == "" {
		message.Message = "Not have system ID."
		return c.JSON(500, message)
	}
	var data struct {
		IntentName string
	}
	if err := c.Bind(&data); err != nil {
		message.Message = err.Error()
		return c.JSON(500, message)
	}
	jwt := string([]rune(authorization)[7:])
	tokens, _ := repositories.DecodeJWT(jwt)
	err := repositories.DeleteIntent(tokens["user_id"].(string), c.QueryParam("systemid"), data.IntentName)
	if err != nil {
		return c.JSON(500, err)
	}
	message.Message = "success"
	return c.JSON(200, message)
}

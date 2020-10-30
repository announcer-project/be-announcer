package handlers

import (
	"be_nms/actions/repositories"

	"github.com/labstack/echo/v4"
)

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

// type Event struct {
// 	ReplyToken string
// 	Type       string
// 	Source     struct {
// 		Type   string
// 		UserID string
// 	}
// 	Message struct {
// 		ID   string
// 		Type string
// 		Text string
// 	}
// }

// type MessageEvent struct {
// 	Destination string
// 	Events      []Event
// }

// func Webhook(c echo.Context) error {
// 	dp := repositories.DialogflowProcessor{}
// 	projectID := "announcer-iysl"
// 	filePath := "dialogflow/announcer-iysl-4dba8d62734e.json"
// 	language := "th"
// 	timeZone := "Asia/Bangkok"

// 	var message struct {
// 		Message string `json:"message"`
// 	}
// 	messageEvent := MessageEvent{}
// 	if err := c.Bind(&messageEvent); err != nil {
// 		message.Message = "server error."
// 		return c.JSON(500, message)
// 	}
// 	log.Print("test ", messageEvent)
// 	dp.Init(projectID, filePath, language, timeZone)
// 	response := dp.ProcessNLP(messageEvent.Events[0].Message.Text, "testUser")
// 	return c.JSON(200, response)
// }

// func CreateIntent(c echo.Context) error {
// 	dp := repositories.DialogflowProcessor{}
// 	projectID := "announcer-iysl"
// 	filePath := "dialogflow/announcer-iysl-4dba8d62734e.json"
// 	language := "th"
// 	timeZone := "Asia/Bangkok"

// 	// var message struct {
// 	// 	Message string `json:"message"`
// 	// }
// 	// messageEvent := MessageEvent{}
// 	// if err := c.Bind(&messageEvent); err != nil {
// 	// 	message.Message = "server error."
// 	// 	return c.JSON(500, message)
// 	// }
// 	// log.Print("test ", messageEvent)
// 	dp.Init(projectID, filePath, language, timeZone)
// 	log.Print("test1")
// 	// response := dp.ProcessNLP(messageEvent.Events[0].Message.Text, "testUser")
// 	response := repositories.CreateIntent(projectID, "หิว", []string{"หิวครับ", "หิวสุด"}, []string{"ไปหาไรกิน", "หาไรกินสิ"})
// 	log.Print(response)
// 	return c.JSON(200, response)
// }

// func ListIntent(c echo.Context) error {
// 	dp := repositories.DialogflowProcessor{}
// 	projectID := "announcer-iysl"
// 	filePath := "dialogflow/announcer-iysl-4dba8d62734e.json"
// 	language := "th"
// 	timeZone := "Asia/Bangkok"

// 	// var message struct {
// 	// 	Message string `json:"message"`
// 	// }
// 	// messageEvent := MessageEvent{}
// 	// if err := c.Bind(&messageEvent); err != nil {
// 	// 	message.Message = "server error."
// 	// 	return c.JSON(500, message)
// 	// }
// 	// log.Print("test ", messageEvent)
// 	dp.Init(projectID, filePath, language, timeZone)
// 	// response := dp.ProcessNLP(messageEvent.Events[0].Message.Text, "testUser")
// 	response, err := repositories.ListIntents(projectID)
// 	if err != nil {
// 		log.Print(err)
// 		return c.JSON(500, err)
// 	}
// 	log.Print(response)
// 	return c.JSON(200, response)
// }

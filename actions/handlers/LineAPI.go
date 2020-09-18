package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"cloud.google.com/go/storage"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

func WebhookLineOA(c echo.Context) error {
	// var projectID, sessionID, languageCode string
	// dialogflow.NewSessionsClient()
	// dialogflow
	explicit("announcer-iysl-4dba8d62734e.json", "announcer-iysl")
	response, err := DetectIntentText("announcer-iysl", "1234567", "สวัสดี", "th-TH")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Output: %s\n", response)
	return c.String(http.StatusOK, "OK!")

}
func DetectIntentText(projectID, sessionID, text, languageCode string) (string, error) {
	ctx := context.Background()

	sessionClient, err := dialogflow.NewSessionsClient(ctx)
	if err != nil {
		return "", err
	}
	defer sessionClient.Close()

	if projectID == "" || sessionID == "" {
		return "", errors.New(fmt.Sprintf("Received empty project (%s) or session (%s)", projectID, sessionID))
	}

	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)
	textInput := dialogflowpb.TextInput{Text: text, LanguageCode: languageCode}
	queryTextInput := dialogflowpb.QueryInput_Text{Text: &textInput}
	queryInput := dialogflowpb.QueryInput{Input: &queryTextInput}
	request := dialogflowpb.DetectIntentRequest{Session: sessionPath, QueryInput: &queryInput}

	response, err := sessionClient.DetectIntent(ctx, &request)
	if err != nil {
		return "", err
	}

	queryResult := response.GetQueryResult()
	fulfillmentText := queryResult.GetFulfillmentText()
	return fulfillmentText, nil
}

func explicit(jsonPath, projectID string) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(jsonPath))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Buckets:")
	it := client.Buckets(ctx, projectID)
	for {
		battrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(battrs.Name)
	}
}

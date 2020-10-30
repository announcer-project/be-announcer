package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

// type DialogflowProcessor struct {
// 	projectID        string
// 	authJSONFilePath string
// 	lang             string
// 	timeZone         string
// 	sessionClient    *dialogflow.SessionsClient
// 	ctx              context.Context
// }

// type NLPResponse struct {
// 	Intent     string            `json:"intent"`
// 	Confidence float32           `json:"confidence"`
// 	Entities   map[string]string `json:"entities"`
// }

// func (dp *DialogflowProcessor) Init(a ...string) (err error) {
// 	dp.projectID = a[0]
// 	dp.authJSONFilePath = a[1]
// 	dp.lang = a[2]
// 	dp.timeZone = a[3]

// 	dp.ctx = context.Background()
// 	sessionClient, err := dialogflow.NewSessionsClient(dp.ctx, option.WithCredentialsFile(dp.authJSONFilePath))
// 	if err != nil {
// 		log.Fatal("Error in auth with Dialogflow")
// 	}
// 	dp.sessionClient = sessionClient

// 	return
// }

// func (dp *DialogflowProcessor) ProcessNLP(rawMessage string, username string) (r NLPResponse) {
// 	sessionID := username
// 	request := dialogflowpb.DetectIntentRequest{
// 		Session: fmt.Sprintf("projects/%s/agent/sessions/%s", dp.projectID, sessionID),
// 		QueryInput: &dialogflowpb.QueryInput{
// 			Input: &dialogflowpb.QueryInput_Text{
// 				Text: &dialogflowpb.TextInput{
// 					Text:         rawMessage,
// 					LanguageCode: dp.lang,
// 				},
// 			},
// 		},
// 		QueryParams: &dialogflowpb.QueryParameters{
// 			TimeZone: dp.timeZone,
// 		},
// 	}
// 	response, err := dp.sessionClient.DetectIntent(dp.ctx, &request)
// 	if err != nil {
// 		log.Fatalf("Error in communication with Dialogflow %s", err.Error())
// 		return
// 	}
// 	queryResult := response.GetQueryResult()
// 	if queryResult.Intent != nil {
// 		r.Intent = queryResult.Intent.DisplayName
// 		r.Confidence = float32(queryResult.IntentDetectionConfidence)
// 	}
// 	r.Entities = make(map[string]string)
// 	params := queryResult.Parameters.GetFields()
// 	if len(params) > 0 {
// 		for paramName, p := range params {
// 			fmt.Printf("Param %s: %s (%s)", paramName, p.GetStringValue(), p.String())
// 			extractedValue := extractDialogflowEntities(p)
// 			r.Entities[paramName] = extractedValue
// 		}
// 	}
// 	return
// }

// func extractDialogflowEntities(p *structpb.Value) (extractedEntity string) {
// 	kind := p.GetKind()
// 	switch kind.(type) {
// 	case *structpb.Value_StringValue:
// 		return p.GetStringValue()
// 	case *structpb.Value_NumberValue:
// 		return strconv.FormatFloat(p.GetNumberValue(), 'f', 6, 64)
// 	case *structpb.Value_BoolValue:
// 		return strconv.FormatBool(p.GetBoolValue())
// 	case *structpb.Value_StructValue:
// 		s := p.GetStructValue()
// 		fields := s.GetFields()
// 		extractedEntity = ""
// 		for key, value := range fields {
// 			if key == "amount" {
// 				extractedEntity = fmt.Sprintf("%s%s", extractedEntity, strconv.FormatFloat(value.GetNumberValue(), 'f', 6, 64))
// 			}
// 			if key == "unit" {
// 				extractedEntity = fmt.Sprintf("%s%s", extractedEntity, value.GetStringValue())
// 			}
// 			if key == "date_time" {
// 				extractedEntity = fmt.Sprintf("%s%s", extractedEntity, value.GetStringValue())
// 			}
// 			// @TODO: Other entity types can be added here
// 		}
// 		return extractedEntity
// 	case *structpb.Value_ListValue:
// 		list := p.GetListValue()
// 		if len(list.GetValues()) > 1 {
// 			// @TODO: Extract more values
// 		}
// 		extractedEntity = extractDialogflowEntities(list.GetValues()[0])
// 		return extractedEntity
// 	default:
// 		return ""
// 	}
// }

func ConnectDialogflow(
	userID,
	systemID,
	projectID,
	authJSONFileBase64,
	lang,
	timeZone string) error {
	db := database.Open()
	defer db.Close()
	system := models.System{}
	db.Where("id = ? and deleted_at is null", systemID).First(&system)
	if system.ID == "" {
		return errors.New("system not found.")
	}
	admin := models.Admin{}
	db.Where("user_id = ? and system_id = ?", userID, system.ID).First(&admin)
	if admin.ID == 0 {
		return errors.New("you not admin.")
	}
	fileByte := Base64toByte(authJSONFileBase64, "json")
	err := ioutil.WriteFile("dialogflow/"+projectID+".json", fileByte, 0700)
	if err != nil {
		return err
	}
	defer os.Remove("dialogflow/" + projectID + ".json")
	sess := ConnectFileStorage()
	if err := CreateFile(sess, fileByte, projectID+".json", "/dialogflow"); err != nil {
		return errors.New("upload auth json file fail.")
	}
	ctx := context.Background()
	_, clientErr := dialogflow.NewIntentsClient(ctx, option.WithCredentialsFile("dialogflow/"+projectID+".json"))
	if clientErr != nil {
		log.Print(clientErr)
		return errors.New(clientErr.Error())
	}

	df := models.DialogflowProcessor{ProjectID: projectID, AuthJSONFilePath: getEnv("STORAGE_PATH", "") + "/dialogflow/" + projectID + ".json", Lang: lang, TimeZone: timeZone}
	system.Dialogflow = df
	db.Save(&system)
	return nil
}

func CreateIntent(projectID, displayName string, trainingPhraseParts, messageTexts []string) error {
	ctx := context.Background()

	intentsClient, clientErr := dialogflow.NewIntentsClient(ctx, option.WithCredentialsFile("dialogflow/announcer-iysl-4dba8d62734e.json"))
	if clientErr != nil {
		log.Print(clientErr)
		return clientErr
	}
	log.Print(intentsClient)
	defer intentsClient.Close()

	if projectID == "" || displayName == "" {
		return errors.New(fmt.Sprintf("Received empty project (%s) or intent (%s)", projectID, displayName))
	}

	parent := fmt.Sprintf("projects/%s/agent", projectID)

	var targetTrainingPhrases []*dialogflowpb.Intent_TrainingPhrase
	var targetTrainingPhraseParts []*dialogflowpb.Intent_TrainingPhrase_Part
	for _, partString := range trainingPhraseParts {
		part := dialogflowpb.Intent_TrainingPhrase_Part{Text: partString}
		targetTrainingPhraseParts = []*dialogflowpb.Intent_TrainingPhrase_Part{&part}
		targetTrainingPhrase := dialogflowpb.Intent_TrainingPhrase{Type: dialogflowpb.Intent_TrainingPhrase_EXAMPLE, Parts: targetTrainingPhraseParts}
		targetTrainingPhrases = append(targetTrainingPhrases, &targetTrainingPhrase)
	}

	intentMessageTexts := dialogflowpb.Intent_Message_Text{Text: messageTexts}
	wrappedIntentMessageTexts := dialogflowpb.Intent_Message_Text_{Text: &intentMessageTexts}
	intentMessage := dialogflowpb.Intent_Message{Message: &wrappedIntentMessageTexts}

	target := dialogflowpb.Intent{DisplayName: displayName, WebhookState: dialogflowpb.Intent_WEBHOOK_STATE_UNSPECIFIED, TrainingPhrases: targetTrainingPhrases, Messages: []*dialogflowpb.Intent_Message{&intentMessage}}

	request := dialogflowpb.CreateIntentRequest{Parent: parent, Intent: &target}

	_, requestErr := intentsClient.CreateIntent(ctx, &request)
	if requestErr != nil {
		return requestErr
	}

	return nil
}

func ListIntents(projectID string) ([]*dialogflowpb.Intent, error) {
	ctx := context.Background()

	intentsClient, clientErr := dialogflow.NewIntentsClient(ctx, option.WithCredentialsFile("dialogflow/announcer-iysl-4dba8d62734e.json"))
	if clientErr != nil {
		return nil, clientErr
	}
	defer intentsClient.Close()

	if projectID == "" {
		return nil, errors.New(fmt.Sprintf("Received empty project (%s)", projectID))
	}

	parent := fmt.Sprintf("projects/%s/agent", projectID)

	request := dialogflowpb.ListIntentsRequest{Parent: parent}

	intentIterator := intentsClient.ListIntents(ctx, &request)
	var intents []*dialogflowpb.Intent

	for intent, status := intentIterator.Next(); status != iterator.Done; {
		intents = append(intents, intent)
		intent, status = intentIterator.Next()
	}

	return intents, nil
}

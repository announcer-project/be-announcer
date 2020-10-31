package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

func CheckConnectDialogflow(userID, systemID string) (bool, error) {
	db := database.Open()
	defer db.Close()
	system := models.System{}
	db.Where("id = ? and deleted_at is null", systemID).First(&system)
	if system.ID == "" {
		return false, errors.New("system not found.")
	}
	admin := models.Admin{}
	db.Where("user_id = ? and system_id = ?", userID, system.ID).First(&admin)
	if admin.ID == 0 {
		return false, errors.New("you not admin.")
	}
	df := models.DialogflowProcessor{}
	db.Where("system_id = ? and deleted_at is null", system.ID).First(&df)
	if df.ID == 0 {
		return false, nil
	}
	return true, nil
}

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

func Webhook(systemid, message string) (interface{}, error) {
	db := database.Open()
	defer db.Close()
	system := models.System{}
	db.Where("id = ? and deleted_at is null", systemid).First(&system)
	if system.ID == "" {
		return nil, errors.New("system not found.")
	}
	df := models.DialogflowProcessor{}
	db.Where("system_id = ? and deleted_at is null", system.ID).First(&df)
	if df.ID == 0 {
		return nil, errors.New("not connect dialogflow.")
	}
	err := DowloadFileJSON(df.AuthJSONFilePath, df.ProjectID+".json")
	if err != nil {
		return nil, err
	}
	defer os.Remove("dialogflow/" + df.ProjectID + ".json")
	df.AuthJSONFilePath = "dialogflow/" + df.ProjectID + ".json"
	err = df.Init()
	if err != nil {
		return nil, err
	}
	response := df.ProcessNLP(message, "testUser")
	return response, nil
}

func DowloadFileJSON(URL, filename string) error {
	response, err := http.Get(URL)
	if err != nil {
	}
	defer response.Body.Close()

	//Create a empty file
	file, err := os.Create("dialogflow/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}

// func CreateIntent(projectID, displayName string, trainingPhraseParts, messageTexts []string) error {

// }

func ListIntents(userID, systemID string) ([]*dialogflowpb.Intent, error) {
	db := database.Open()
	defer db.Close()
	system := models.System{}
	db.Where("id = ? and deleted_at is null", systemID).First(&system)
	if system.ID == "" {
		return nil, errors.New("system not found.")
	}
	admin := models.Admin{}
	db.Where("user_id = ? and system_id = ?", userID, system.ID).First(&admin)
	if admin.ID == 0 {
		return nil, errors.New("you not admin.")
	}
	df := models.DialogflowProcessor{}
	db.Where("system_id = ? and deleted_at is null", system.ID).First(&df)
	if df.ID == 0 {
		return nil, errors.New("system not connect dialogflow.")
	}
	err := DowloadFileJSON(df.AuthJSONFilePath, df.ProjectID+".json")
	if err != nil {
		return nil, err
	}
	defer os.Remove("dialogflow/" + df.ProjectID + ".json")
	df.AuthJSONFilePath = "dialogflow/" + df.ProjectID + ".json"
	err = df.Init()
	ctx := context.Background()

	intentsClient, clientErr := dialogflow.NewIntentsClient(ctx, option.WithCredentialsFile(df.AuthJSONFilePath))
	if clientErr != nil {
		return nil, clientErr
	}
	defer intentsClient.Close()

	parent := fmt.Sprintf("projects/%s/agent", df.ProjectID)

	request := dialogflowpb.ListIntentsRequest{Parent: parent}

	intentIterator := intentsClient.ListIntents(ctx, &request)
	var intents []*dialogflowpb.Intent

	for intent, status := intentIterator.Next(); status != iterator.Done; {
		intents = append(intents, intent)
		intent, status = intentIterator.Next()
	}

	return intents, nil
}

func GetIntent(userID, systemID, IntentName string) (*dialogflowpb.Intent, error) {
	db := database.Open()
	defer db.Close()
	system := models.System{}
	db.Where("id = ? and deleted_at is null", systemID).First(&system)
	if system.ID == "" {
		return nil, errors.New("system not found.")
	}
	admin := models.Admin{}
	db.Where("user_id = ? and system_id = ?", userID, system.ID).First(&admin)
	if admin.ID == 0 {
		return nil, errors.New("you not admin.")
	}
	df := models.DialogflowProcessor{}
	db.Where("system_id = ? and deleted_at is null", system.ID).First(&df)
	if df.ID == 0 {
		return nil, errors.New("system not connect dialogflow.")
	}
	err := DowloadFileJSON(df.AuthJSONFilePath, df.ProjectID+".json")
	if err != nil {
		return nil, err
	}
	defer os.Remove("dialogflow/" + df.ProjectID + ".json")
	df.AuthJSONFilePath = "dialogflow/" + df.ProjectID + ".json"
	err = df.Init()

	ctx := context.Background()
	c, err := dialogflow.NewIntentsClient(ctx, option.WithCredentialsFile(df.AuthJSONFilePath))
	if err != nil {
		return nil, err
	}
	defer c.Close()
	// name:="projects/<Project ID>/agent/intents/<Intent ID>"
	req := &dialogflowpb.GetIntentRequest{Name: IntentName}
	response, err := c.GetIntent(ctx, req)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func CreateIntent(userID, systemID, displayName string, trainingPhraseParts, messageTexts []string) error {
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
	df := models.DialogflowProcessor{}
	db.Where("system_id = ? and deleted_at is null", system.ID).First(&df)
	if df.ID == 0 {
		return errors.New("system not connect dialogflow.")
	}
	err := DowloadFileJSON(df.AuthJSONFilePath, df.ProjectID+".json")
	if err != nil {
		return err
	}
	defer os.Remove("dialogflow/" + df.ProjectID + ".json")
	df.AuthJSONFilePath = "dialogflow/" + df.ProjectID + ".json"
	err = df.Init()

	ctx := context.Background()

	intentsClient, clientErr := dialogflow.NewIntentsClient(ctx, option.WithCredentialsFile(df.AuthJSONFilePath))
	if clientErr != nil {
		log.Print(clientErr)
		return clientErr
	}
	defer intentsClient.Close()

	parent := fmt.Sprintf("projects/%s/agent", df.ProjectID)

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

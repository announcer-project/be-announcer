package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func GetUserBySocialId(SocialID, Social string) (interface{}, error) {
	db := database.Open()
	user := models.User{}
	if Social == "line" {
		db.First(&user, "line_id = ?", SocialID)
		if user.ID == "" {
			return nil, errors.New("You don't register.")
		}
	} else if Social == "facebook" {
		db.First(&user, "facebook_id = ?", SocialID)
		if user.ID == "" {
			return nil, errors.New("You don't register.")
		}
	}
	return user, nil
}

type UserLine struct {
	UserId     string `json: "user_id"`
	PictureUrl string `json: "picture_url"`
	Email      string `json: "email"`
}

func LineLogin(code string) (interface{}, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", getEnv("CLIENT_ID", ""))
	data.Set("client_secret", getEnv("CLIENT_SECRET", ""))
	data.Set("code", code)
	data.Set("redirect_uri", getEnv("REDIRECT_URI", ""))

	client := &http.Client{}

	request, err := http.NewRequest("POST", "https://api.line.me/oauth2/v2.1/token", strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		log.Fatalln(err)
	}

	res, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var token struct {
		Access_token string `json: "token"`
		ID_token     string `json:"id_token"`
	}

	json.Unmarshal(body, &token)
	log.Print("token ", token)
	tokens, _ := DecodeJWT(token.ID_token)
	log.Print("email ", tokens["email"])
	userline, _ := GetUserProfileLine(token.Access_token, tokens["email"].(string))
	user, err := GetUserBySocialId(userline.(UserLine).UserId, "line")
	if err != nil {
		return userline, err
	}
	return user, nil
}

func GetUserProfileLine(AccessToken string, email string) (interface{}, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.line.me/v2/profile", nil)
	req.Header.Set("Authorization", "Bearer "+AccessToken)
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	log.Print("body", body)
	user := UserLine{}
	json.Unmarshal(body, &user)
	user.Email = email
	return user, nil
}

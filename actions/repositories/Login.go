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

	"github.com/labstack/echo/v4"
)

type Token struct {
	ID_Token string `json:"id_token"`
}

type UserIdSocial struct {
	UserId string `json: "userId"`
}

func GetUserIDLine(c echo.Context) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", getEnv("CLIENT_ID", ""))
	data.Set("client_secret", getEnv("CLIENT_SECRET", ""))
	data.Set("code", c.FormValue("code"))
	data.Set("redirect_uri", getEnv("REDIRECT_URI", ""))
	client := &http.Client{}
	request, _ := http.NewRequest("POST", "https://api.line.me/oauth2/v2.1/token", strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(request)
	defer res.Body.Close()
	if err != nil {
		return "nil", err
	}
	body, _ := ioutil.ReadAll(res.Body)
	log.Print(string(body))
	token := Token{}
	json.Unmarshal(body, &token)
	profile, _ := DecodeJWT(token.ID_Token)
	return profile["sub"].(string), nil
}

func GetUserBySocialId(UserId, Social string) (interface{}, error) {
	db := database.Open()
	user := models.User{}
	if Social == "line" {
		db.First(&user, "line_id = ?", UserId)
		if user.ID == "" {
			return nil, errors.New("You don't register.")
		}
	}
	return user, nil
}

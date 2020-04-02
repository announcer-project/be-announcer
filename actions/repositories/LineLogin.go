package repositories

import (
	"be_nms/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetAccessTokenLine(c echo.Context) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", getEnv("CLIENT_ID", ""))
	data.Set("client_secret", getEnv("CLIENT_SECRET", ""))
	data.Set("code", c.FormValue("code"))
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
	token := models.Token{}
	json.Unmarshal(body, &token)

	return token.Access_token, nil
}

func GetUserProfileLine(AccessToken string) (interface{}, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.line.me/v2/profile", nil)
	req.Header.Set("Authorization", "Bearer "+AccessToken)
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	user := models.LineProfile{}
	json.Unmarshal(body, &user)
	return user, nil
}

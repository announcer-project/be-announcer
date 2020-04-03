package repositories

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

type Token struct {
	Access_token string `json: "token"`
}

type UserIdSocial struct {
	UserId string `json: "userId"`
}

func GetAccessTokenLine(c echo.Context) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", getEnv("CLIENT_ID", ""))
	data.Set("client_secret", getEnv("CLIENT_SECRET", ""))
	data.Set("code", c.FormValue("code"))
	data.Set("redirect_uri", getEnv("REDIRECT_URI", ""))
	client := &http.Client{}
	request, _ := http.NewRequest("POST", "https://api.line.me/oauth2/v2.1/token", strings.NewReader(data.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, _ := client.Do(request)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	token := Token{}
	json.Unmarshal(body, &token)
	return token.Access_token, nil
}

func GetUserIdLine(AccessToken string) (interface{}, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.line.me/v2/profile", nil)
	req.Header.Set("Authorization", "Bearer "+AccessToken)
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	user := UserIdSocial{}
	json.Unmarshal(body, &user)
	return user.UserId, nil
}

func GetUserBySocialId(UserId, Social string) (bool, err) {
	if Social == "line" {
		//check with db
		//user
	}
	return true, nil
}

package repositories

import (
	"be_nms/database"
	"be_nms/models"
	"errors"
)

type Token struct {
	ID_Token string `json:"id_token"`
}

type UserIdSocial struct {
	UserId string `json: "userId"`
}

// func GetUserIDLine(c echo.Context) (string, error) {
// 	code := c.Request().Header.Get("Code")
// 	data := url.Values{}
// 	data.Set("grant_type", "authorization_code")
// 	data.Set("client_id", getEnv("CLIENT_ID", ""))
// 	data.Set("client_secret", getEnv("CLIENT_SECRET", ""))
// 	data.Set("code", code)
// 	data.Set("redirect_uri", getEnv("REDIRECT_URI", ""))
// 	client := &http.Client{}
// 	request, _ := http.NewRequest("POST", "https://api.line.me/oauth2/v2.1/token", strings.NewReader(data.Encode()))
// 	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 	res, err := client.Do(request)
// 	defer res.Body.Close()
// 	if err != nil {
// 		return "nil", err
// 	}
// 	body, _ := ioutil.ReadAll(res.Body)
// 	token := Token{}
// 	json.Unmarshal(body, &token)
// 	profile, _ := DecodeJWT(token.ID_Token)
// 	return profile["sub"].(string), nil
// }

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

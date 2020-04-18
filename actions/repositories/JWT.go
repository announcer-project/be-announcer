package repositories

import (
	"be_nms/models"
	"log"

	jwtgo "github.com/dgrijalva/jwt-go"
)

type JWT struct {
	UserID string `json:"user_id"`
	FName  string `json:"fname"`
	LName  string `json:"lname"`
	jwtgo.StandardClaims
}

func EncodeJWT(user models.User) string {
	tokenJWT := JWT{UserID: user.ID, FName: user.FName, LName: user.LName}
	token := jwtgo.NewWithClaims(jwtgo.GetSigningMethod("HS256"), tokenJWT)
	jwt, err := token.SignedString([]byte("newsmanagement"))
	if err != nil {
		log.Fatalln(err)
	}
	return jwt
}

func DecodeJWT(jwt string) (jwtgo.MapClaims, interface{}) {
	token, _ := jwtgo.Parse(jwt, func(token *jwtgo.Token) (interface{}, error) {
		return []byte("newsmanagement"), nil
	})
	tokens := token.Claims.(jwtgo.MapClaims)
	return tokens, nil
}

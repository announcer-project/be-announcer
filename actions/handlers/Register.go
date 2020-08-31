package handlers

import (
	"be_nms/actions/repositories"
	"be_nms/models"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	user, err := repositories.Register(c.FormValue("email"), c.FormValue("fname"), c.FormValue("lname"), c.FormValue("line"), c.FormValue("facebook"), c.FormValue("google"), c.FormValue("imagesocial"), c.FormValue("imageUrl"), c.FormValue("imageProfile"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	jwt := repositories.EncodeJWT(user.(models.User))
	return c.JSON(http.StatusOK, jwt)
}

func SendOTP(c echo.Context) error {
	otp := c.FormValue("otp")
	email := c.FormValue("email")
	go repositories.SendEmail("OTP", email, otp)
	return c.JSON(http.StatusOK, "Send OTP Success!")
}

func CheckUserByLineID(c echo.Context) error {
	lineid := c.FormValue("lineid")
	haveuser := repositories.CheckUserByLineID(lineid)
	if haveuser {
		return c.JSON(http.StatusOK, haveuser)
	}
	return c.JSON(http.StatusOK, haveuser)
}

func CheckUserByEmail(c echo.Context) error {
	email := c.FormValue("email")
	user, err := repositories.CheckUserByEmail(email)
	if err != nil {
		log.Print("have account")
		return c.JSON(http.StatusBadRequest, user)
	}
	log.Print("not have account")
	return c.JSON(http.StatusOK, "Not have account")
}

func ConnectSocialWithAccount(c echo.Context) error {
	social := c.FormValue("social")
	socialID := c.FormValue("socialid")
	userID := c.FormValue("userid")
	log.Print("connect ", social, socialID, userID)
	user, err := repositories.ConnectSocialWithAccount(social, socialID, userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	jwt := repositories.EncodeJWT(user.(models.User))
	return c.JSON(http.StatusOK, jwt)
}

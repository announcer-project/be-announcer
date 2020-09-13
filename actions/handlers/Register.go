package handlers

import (
	"be_nms/actions/repositories"
	"be_nms/models"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Register(c echo.Context) error {
	var data struct {
		Email        string
		Fname        string
		Lname        string
		Line         string
		Facebook     string
		Google       string
		ImageSocial  bool
		ImageUrl     string
		ImageProfile string
	}
	if err := c.Bind(&data); err != nil {
		log.Print("error ", err)
		return err
	}
	user, err := repositories.Register(
		data.Email,
		data.Fname,
		data.Lname,
		data.Line,
		data.Facebook,
		data.Google,
		data.ImageSocial,
		data.ImageUrl,
		data.ImageProfile)
	if err != nil {
		fail := struct {
			Message string `json:"message"`
		}{
			err.Error(),
		}
		if fail.Message == "Register fail." {
			return c.JSON(500, fail)
		} else {
			return c.JSON(400, fail)
		}
	}
	jwt := repositories.EncodeJWT(user.(models.User))
	success := struct {
		JWT string `json:"jwt"`
	}{
		jwt,
	}
	return c.JSON(http.StatusOK, success)
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

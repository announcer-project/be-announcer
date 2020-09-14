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
	var data struct {
		OTP   string
		Email string
	}
	if err := c.Bind(&data); err != nil {
		log.Print("error ", err)
		return err
	}
	go repositories.SendEmail("OTP", data.Email, data.OTP)
	success := struct {
		Message string `json:"message"`
	}{
		"send otp success.",
	}
	return c.JSON(http.StatusOK, success)
}

func CheckUserByLineID(c echo.Context) error {
	var data struct {
		LineID string
	}
	if err := c.Bind(&data); err != nil {
		log.Print("error ", err)
		return err
	}
	haveuser := repositories.CheckUserByLineID(data.LineID)
	var message struct {
		Message string `json:"message"`
	}
	if haveuser {
		message.Message = "have account."
		return c.JSON(http.StatusOK, message)
	}
	message.Message = "not have account."
	return c.JSON(http.StatusOK, message)
}

func CheckUserByEmail(c echo.Context) error {
	var data struct {
		Email string
	}
	if err := c.Bind(&data); err != nil {
		log.Print("error ", err)
		return err
	}
	_, err := repositories.CheckUserByEmail(data.Email)
	var message struct {
		Message string `json:"message"`
	}
	if err != nil {
		message.Message = "have account."
		return c.JSON(400, message)
	}
	message.Message = "not have account."
	return c.JSON(http.StatusOK, message)
}

func ConnectSocialWithAccount(c echo.Context) error {
	var data struct {
		Social   string
		SocialID string
		UserID   string
	}
	if err := c.Bind(&data); err != nil {
		log.Print("error ", err)
		return err
	}
	user, err := repositories.ConnectSocialWithAccount(data.Social, data.SocialID, data.UserID)
	if err != nil {
		message := struct {
			Message string `json:"message"`
		}{
			err.Error(),
		}
		return c.JSON(http.StatusBadRequest, message)
	}
	jwt := repositories.EncodeJWT(user.(models.User))
	success := struct {
		JWT string `json:"jwt"`
	}{
		jwt,
	}
	return c.JSON(http.StatusOK, success)
}

package handlers

import (
	"be_nms/actions/repositories"
	"be_nms/models"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	var data struct {
		Social   string
		SocialID string
	}
	if err := c.Bind(&data); err != nil {
		log.Print("error ", err)
		return err
	}
	user, err := repositories.GetUserBySocialId(data.SocialID, data.Social)
	if err != nil {
		fail := struct {
			Message  string `json:"message"`
			SocialID string `json:"social_id"`
		}{
			err.Error(),
			data.SocialID,
		}
		return c.JSON(401, fail)
	}
	jwt := repositories.EncodeJWT(user.(models.User))
	success := struct {
		JWT string `json:"jwt"`
	}{
		jwt,
	}
	return c.JSON(http.StatusOK, success)
}

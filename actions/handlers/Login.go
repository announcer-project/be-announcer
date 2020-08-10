package handlers

import (
	"be_nms/actions/repositories"
	"be_nms/models"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	social := c.Request().Header.Get("Social")
	socialID := c.Request().Header.Get("SocialID")
	email := c.Request().Header.Get("Email")
	pictureUrl := c.Request().Header.Get("PictureUrl")
	log.Print(social, socialID, email, pictureUrl)
	// userID := ""
	// err := errors.New("error")
	// if social == "line" {
	// 	userID, err = repositories.GetUserIDLine(c)
	// 	if err != nil {
	// 		return c.JSON(400, err)
	// 	}
	// } else if social == "facebook" {
	// 	userID = c.Request().Header.Get("UserID")
	// }
	user, err := repositories.GetUserBySocialId(socialID, social)
	if err != nil {
		return c.JSON(http.StatusBadRequest, socialID)
	}
	jwt := repositories.EncodeJWT(user.(models.User))
	return c.JSON(http.StatusOK, jwt)
}

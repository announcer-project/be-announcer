package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func WebhookLineOA(c echo.Context) error {
	return c.String(http.StatusOK, "OK!")
}

// func SetDefaultRichMenuRegister(c echo.Context) error {
// 	richmenuid, err := repositories.CreateRichmenu(c)
// 	if err != nil {
// 		return err
// 	}
// 	err = repositories.SetImageToRichMenu(c, richmenuid.(string))
// 	if err != nil {
// 		return err
// 	}
// 	err = repositories.SetDefaultRichMenu(c, richmenuid.(string))
// 	if err != nil {
// 		return err
// 	}
// 	return c.JSON(http.StatusOK, richmenuid)
// }

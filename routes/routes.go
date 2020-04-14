package routes

import (
	"be_nms/actions/handlers"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func Init() *echo.Echo {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to News Management System")
	})
	e.POST("/linelogin", handlers.LineLogin)
	e.POST("/register", handlers.Register)

	e.POST("/createnews", handlers.CreateNews)
	e.GET("/getnews/:id", handlers.GetNewsByID)
	return e
}

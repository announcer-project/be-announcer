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
	//Account
	e.POST("/linelogin", handlers.LineLogin)
	e.POST("/register", handlers.Register)
	//System
	e.GET("/getallsystems", handlers.GetAllSystems)
	//NewsManagement
	e.GET("/news/getallnews", handlers.GetAllNews)
	e.POST("/createnews", handlers.CreateNews)
	e.GET("/news/:id", handlers.GetNewsByID)
	e.POST("/announcenews", handlers.AnnounceNews)
	e.POST("/createtargetgroup", handlers.CreateTargetGroup)
	//Social
	e.POST("/webhooklineoa", handlers.WebhookLineOA)
	return e
}

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
	e.GET("/system/allsystem", handlers.GetAllSystems)
	e.POST("/system/create", handlers.CreateSystem)
	//NewsManagement
	e.GET("/news/allnews", handlers.GetAllNews)
	e.POST("/news/create", handlers.CreateNews)
	e.GET("/news/:id", handlers.GetNewsByID)
	e.POST("/news/newstype/create", handlers.CreateNewsType)
	e.GET("/news/newstype/allnewstype", handlers.GetAlNewsType)
	e.POST("/announcenews", handlers.AnnounceNews)
	//TargetGroup
	e.POST("/targetgroup/create", handlers.CreateTargetGroup)
	e.GET("/targetgroup/all", handlers.GetAllTargetGroup)
	//Social
	e.POST("/webhooklineoa", handlers.WebhookLineOA)
	return e
}

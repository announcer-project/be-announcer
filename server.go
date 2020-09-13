package main

import (
	"be_nms/database"
	"be_nms/routes"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4/middleware"
)

func getPort() string {
	var port = os.Getenv("PORT") // ----> (A)
	if port == "" {
		port = "8080"
		fmt.Println("No Port In Heroku" + port)
	}
	return ":" + port // ----> (B)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	e := routes.Init()
	db := database.Open()
	defer db.Close()
	// database.Migration(db)
	// database.SetData(db)
	e.Use(middleware.CORS())
	e.Logger.Fatal(e.Start(getPort()))
}

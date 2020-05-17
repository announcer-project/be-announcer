package main

import (
	"be_nms/database"
	"be_nms/routes"
	"fmt"
	"os"

	_ "github.com/jinzhu/gorm/dialects/mysql"
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
	e := routes.Init()
	db := database.Open()
	defer db.Close()
	// database.Migration(db)
	// database.SetData(db)
	e.Use(middleware.CORS())
	e.Logger.Fatal(e.Start(getPort()))
}

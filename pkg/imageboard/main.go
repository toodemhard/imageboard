package imageboard

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Run() {
	db := ConnectToDB()
	_ = db

	e := echo.New()
	e.Static("/", "static")
	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

package imageboard

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func postThread(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		thread := Thread{}
		if err := c.Bind(&thread); err != nil {
			return err
		}
		fmt.Println(thread)
		CreateThread(db, thread)
		return c.HTML(http.StatusOK, `<div>done it m8<\div>`)

	}
}

func Run() {
	db := ConnectToDB()

	e := echo.New()
	e.Static("/", "static")
	e.POST("/thread", postThread(db))
	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

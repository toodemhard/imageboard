package imageboard

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const assetsDir = "public/assets/"

func Run() {
	db := ConnectToDB()
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.HideBanner = true

	template, err := ParseTemplates()
	if err != nil {
		panic(err)
	}
	e.Renderer = &StaticTemplate{template}

	if len(os.Args) > 1 {
		flag := os.Args[1]
		if flag == "--dev" || flag == "-d" {
			e.Renderer = &LiveTemplate{}
		}
	}

	h := Handler{db}

	e.GET("/images/:id", getImage)
	e.Static("/assets", assetsDir)
	e.POST("/thread", h.postThread)
	e.GET("/thread/:id", h.getThread)
	e.POST("/thread/:id/reply", h.postReply)
	e.GET("/", h.getIndex)
	if err := e.Start(":8080"); err != http.ErrServerClosed {
		e.Logger.Fatal(err)
	}
}

package imageboard

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	txt "text/template"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const publicDir = "../../public/"

func Run() {
	if len(os.Args) > 1 {
		flag := os.Args[1]
		if flag == "--dev" || flag == "-d" {
			fmt.Println("dev mode")
		}
	}

	db := ConnectToDB()
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.HideBanner = true
	e.Renderer = &Template{
		template.Must(template.ParseGlob(publicDir + "pages/*.html")),
		txt.Must(txt.ParseFiles(publicDir + "layout.html")),
	}

	h := Handler{db}

	e.GET("/images/:id", getImage)
	e.Static("/assets", publicDir+"assets/")
	e.POST("/thread", h.postThread)
	e.GET("/thread/:id", h.getThread)
	e.POST("/thread/:id/reply", h.postReply)
	e.GET("/", h.getIndex)
	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

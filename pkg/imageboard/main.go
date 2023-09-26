package imageboard

import (
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

const publicDir = "../../public/"

func postThread(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		thread := Thread{}
		if err := c.Bind(&thread); err != nil {
			return err
		}
		CreateThread(db, thread)
		return c.HTML(http.StatusOK, `<div>done it m8`)
	}
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func getThreads(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		threads, err := GetAllThreads(db)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "db query failed")
		}
		return c.Render(http.StatusOK, "thread.html", threads)
	}
}

func Run() {
	db := ConnectToDB()
	e := echo.New()

	t := &Template{
		template.Must(template.ParseGlob(publicDir + "*.html")),
	}

	e.Renderer = t

	e.Static("/assets", publicDir+"assets")
	e.POST("/thread", postThread(db))
	e.GET("/", getThreads(db))
	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

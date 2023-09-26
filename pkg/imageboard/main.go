package imageboard

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"

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

func getThread(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "invalid thread id")
		}

		thread, err := queryThread(db, int64(id))
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "db query failed")
		}

		replies, err := queryThreadReplies(db, id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "db query failed")
		}

		page := struct {
			Thread  Thread
			Replies []Reply
		}{
			thread,
			replies,
		}

		return c.Render(http.StatusOK, "thread.html", page)
	}
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	err := t.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		log.Println(err)
	}
	return err
}

func index(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		threads, err := queryAllThreads(db)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "db query failed")
		}
		return c.Render(http.StatusOK, "index.html", threads)
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
	e.GET("/thread/:id", getThread(db))
	e.GET("/", index(db))
	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

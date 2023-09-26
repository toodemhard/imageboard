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

func (h *Handler) postThread(c echo.Context) error {
	thread := Thread{}
	if err := c.Bind(&thread); err != nil {
		return err
	}
	CreateThread(h.db, thread)
	return c.HTML(http.StatusOK, `<div>done it m8`)
}

func (h *Handler) getThread(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "invalid thread id")
	}

	thread, err := queryThread(h.db, int64(id))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "db query failed")
	}

	replies, err := queryThreadReplies(h.db, id)
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

func (h *Handler) index(c echo.Context) error {
	threads, err := queryAllThreads(h.db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "db query failed")
	}
	return c.Render(http.StatusOK, "index.html", threads)
}

type Handler struct {
	db *sqlx.DB
}

func Run() {
	db := ConnectToDB()
	e := echo.New()

	e.Renderer = &Template{
		template.Must(template.ParseGlob(publicDir + "*.html")),
	}

	h := Handler{db}

	e.Static("/assets", publicDir+"assets")
	e.POST("/thread", h.postThread)
	e.GET("/thread/:id", h.getThread)
	// e.POST("/thread/:id/reply", getThread(db))
	e.GET("/", h.index)
	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

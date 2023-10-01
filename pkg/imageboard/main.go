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

type Handler struct {
	db *sqlx.DB
}

func (h *Handler) postReply(c echo.Context) error {
	thread_id, _ := strconv.Atoi(c.Param("id"))
	img, _ := c.FormFile("image")
	reply := Reply{}
	if err := c.Bind(&reply); err != nil {
		return err
	}

	reply.Thread_id = thread_id
	CreateReply(h.db, reply, img)

	return c.HTML(http.StatusOK, "<div>submitted")
}

func (h *Handler) postThread(c echo.Context) error {
	thread := Thread{}
	img, err := c.FormFile("image")
	if err != nil {
		log.Println(err)
		return err
	}

	if err := c.Bind(&thread); err != nil {
		return err
	}

	if err := CreateThread(h.db, thread, img); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

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
		log.Println(err)
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

func (h *Handler) getIndex(c echo.Context) error {
	threads, err := queryAllThreads(h.db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.Render(http.StatusOK, "index.html", threads)
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

func getImage(c echo.Context) error {
	return c.File(imagesDir + c.Param("id"))
}

func Run() {
	db := ConnectToDB()
	e := echo.New()

	e.Renderer = &Template{
		template.Must(template.ParseGlob(publicDir + "*.html")),
	}

	h := Handler{db}

	e.GET("/images/:id", getImage)
	e.Static("/assets", publicDir+"assets")
	e.POST("/thread", h.postThread)
	e.GET("/thread/:id", h.getThread)
	e.POST("/thread/:id/reply", h.postReply)
	e.GET("/", h.getIndex)
	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

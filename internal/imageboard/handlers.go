package imageboard

import (
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

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
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound)
	}

	thread, err := queryThread(h.db, int64(id))
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound)
	}

	replies, err := queryThreadReplies(h.db, id)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound)
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
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.Render(http.StatusOK, "index.html", threads)
}

func getImage(c echo.Context) error {
	return c.File(imagesDir + c.Param("id"))
}

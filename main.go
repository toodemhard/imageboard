package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type post struct {
	title       string
	description string
	id          int
	replies     []reply
}

type reply struct {
	text string
}

var x int
var posts []post

func getPosts(c echo.Context) error {
	string := ""
	for _, p := range posts {
		string += fmt.Sprintf("%v\n", p)
	}
	return c.String(http.StatusOK, string)
}

func getComponent(c echo.Context) error {
	component := "<h1>yo bitsh</h1>"
	return c.HTML(http.StatusOK, component)
}

func main() {

	e := echo.New()

	e.Static("/", "static")
	e.GET("/component", getComponent)
	e.Start(":8080")
}

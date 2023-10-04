package imageboard

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type StaticTemplate struct {
	templates *template.Template
}

type LiveTemplate struct {
}

type Page struct {
	Content string
}

func ParseTemplates() (*template.Template, error) {
	t, err := template.ParseGlob(publicDir + "views/*.html")
	if err != nil {
		return nil, err
	}

	t.ParseGlob(publicDir + "components/*.html")
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *StaticTemplate) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	err := t.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	return nil
}

func (t *LiveTemplate) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	templates, err := ParseTemplates()
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	err = templates.ExecuteTemplate(w, name, data)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	return nil
}

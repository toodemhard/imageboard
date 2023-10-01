package imageboard

import (
	"bytes"
	"html/template"
	"io"
	txt "text/template"

	"github.com/labstack/echo/v4"
)

type Template struct {
	pages  *template.Template
	layout *txt.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	var pageContent bytes.Buffer
	err := t.pages.ExecuteTemplate(&pageContent, name, data)
	if err != nil {
		return err
	}

	page := struct {
		Content string
	}{
		pageContent.String(),
	}

	t.layout.Execute(w, page)
	return err
}

package imageboard

import (
	"html/template"
	"io"
	"os"
	txt "text/template"

	"github.com/labstack/echo/v4"
)

const templatesDir = "templates/"
const srcDir = templatesDir + "src/"
const outDir = templatesDir + "out/"

type StaticTemplate struct {
	pages *template.Template
}

type LiveTemplate struct {
}

type Page struct {
	Content string
}

func ParseTemplates() (*template.Template, error) {
	err := ComposeTemplates()
	if err != nil {
		return nil, err
	}
	template, err := template.ParseGlob(outDir + "*.html")
	if err != nil {
		return nil, err
	}
	return template, nil
}

func (t *StaticTemplate) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	err := t.pages.ExecuteTemplate(w, name, data)
	if err != nil {
		return err
	}

	return nil
}

func (l *LiveTemplate) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	pages, err := ParseTemplates()
	if err != nil {
		c.Logger().Error(err)
		return err
	}
	err = pages.ExecuteTemplate(w, name, data)
	if err != nil {
		return err
	}

	return nil
}

func ComposeTemplates() error {
	base := txt.Must(txt.ParseFiles(srcDir + "base.html"))

	pages, err := os.ReadDir(srcDir + "pages")
	if err != nil {
		return err
	}

	for _, p := range pages {
		f, err := os.Create(outDir + p.Name())
		defer f.Close()
		if err != nil {
			return err
		}

		content, err := os.ReadFile(srcDir + "pages/" + p.Name())
		page := Page{string(content)}
		if err != nil {
			return err
		}

		err = base.Execute(f, page)
		if err != nil {
			return err
		}
	}

	return nil
}

// type LiveTemplate struct {
// }
//
// func (l *LiveTemplate) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
//
// }

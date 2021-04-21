package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
)

type Template struct {
	templates *template.Template
}

func NewTemplate() *Template {
	tt, err := template.ParseFS(public, "public/*.gohtml")
	if err != nil {
		panic(err)
	}
	return &Template{templates: tt}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, fmt.Sprintf("%s.gohtml", name), data)
}

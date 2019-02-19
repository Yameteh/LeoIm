package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	//	t := &Template{
	//		templates: template.Must(template.ParseGlob("*.html")),
	//	}

	e := echo.New()
	//e.Renderer = t
	e.Static("/", "assets")
	e.GET("/hello", Hello)
	e.GET("/dologin", Login)
	e.Start(":8081")

}

func Hello(c echo.Context) error {
	id := c.QueryParam("id")
	return c.String(http.StatusOK, string(id))
}

func Login(c echo.Context) error {
	fmt.Println("login")
	acc := c.QueryParam("account")
	pwd := c.QueryParam("password")
	rem := c.QueryParam("remember")
	fmt.Printf("acc %s pwd %s rem %s", acc, pwd, rem)
	return c.String(http.StatusOK, "login ")
}

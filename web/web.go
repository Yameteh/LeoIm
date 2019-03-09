package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"
)

const (
	WEB_INI = "web.ini"
)

var config *Config
var storeManager *StoreManager
var sessionManager *SessionManager

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	config = NewConfig(WEB_INI)
	if config != nil {
		config.Parse()
	} else {
		glog.Error("router config ini missed")
		return
	}
	storeManager = NewStoreManager()
	err := storeManager.Init()
	if err != nil {
		glog.Error(err)
	}
	sessionManager = NewSessionManager()
	go backendListen()
	frontListen()
}

func frontListen() {
	e := echo.New()
	t := &Template{
		templates: template.Must(template.ParseGlob("assets/*.html")),
	}
	e.Renderer = t
	e.Static("/", "static")
	e.GET("/login", Login)
	e.GET("/dologin", DoLogin)
	e.GET("/loginout", LoginOut)
	e.GET("/index", Index)
	e.GET("/account", Account)
	e.GET("/message", Message)
	e.GET("/console", Console)
	e.GET("/about", About)
	e.Use(session.Middleware(sessionManager.GetStore()))
	e.Start(fmt.Sprintf("%s:%d", config.Domain, config.Port+1))
}

func Account(c echo.Context) error {
	if !sessionManager.IsNewSession(c) {
		users, err := storeManager.GetOnlineUsers(0, 0)
		fmt.Println(len(users))
		if err != nil {
			glog.Error(err)
		}
		return c.Render(http.StatusOK, "account", users)
	} else {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
}

func Message(c echo.Context) error {
	if !sessionManager.IsNewSession(c) {
		return c.Render(http.StatusOK, "message", "")
	} else {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
}

func Console(c echo.Context) error {
	if !sessionManager.IsNewSession(c) {
		return c.Render(http.StatusOK, "console", "")
	} else {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
}

func About(c echo.Context) error {
	if !sessionManager.IsNewSession(c) {
		return c.Render(http.StatusOK, "about", "")
	} else {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
}

func Index(c echo.Context) error {
	if !sessionManager.IsNewSession(c) {
		return c.Render(http.StatusOK, "index", config.FrontTitle)
	} else {
		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
}

func Login(c echo.Context) error {
	if !sessionManager.IsNewSession(c) {
		return c.Redirect(http.StatusTemporaryRedirect, "/index")
	} else {
		return c.Render(http.StatusOK, "login", config.FrontTitle)
	}
}

func LoginOut(c echo.Context) error {
	sessionManager.Remove(c)
	return c.Redirect(http.StatusTemporaryRedirect, "/login")
}

func DoLogin(c echo.Context) error {
	fmt.Println("login")
	acc := c.QueryParam("account")
	pwd := c.QueryParam("password")
	rem := c.QueryParam("remember")
	fmt.Printf("acc %s pwd %s rem %s", acc, pwd, rem)
	has := storeManager.CheckAdmin(acc, pwd)
	if has {
		fmt.Println("login success")
		sessionManager.Save(c)
		c.Response().Header().Add("Cache-Control", "no-store")
		return c.Redirect(http.StatusPermanentRedirect, "/index")
	} else {
		fmt.Println("login failed")
		return c.Redirect(http.StatusOK, "/login")
	}
}

func backendListen() {
	e := echo.New()
	e.Use(middleware.BasicAuth(basicAuthHandle))
	e.Static("file", config.RootDirectory)
	AddRouter(e, &FileRouter{})
	e.Start(fmt.Sprintf("%s:%d", config.Domain, config.Port))
}

func basicAuthHandle(u, p string, c echo.Context) (bool, error) {
	user := RedisQueryUser(u)
	if user != nil && user.Uuid == u && user.Token == p {
		fmt.Println("basic auth true")
		return true, nil
	} else {
		fmt.Println("base auth false")
		return false, nil
	}
}

package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	WEB_INI = "web.ini"
)

var config *Config

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

package main

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/labstack/echo"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func AddRouter(e *echo.Echo, r Router) {
	e.Match(r.Methods(), r.Path(), r.Handler)
}

type Router interface {
	Methods() []string
	Path() string
	Handler(c echo.Context) error
}

type FileRouter struct {
}

func (fr *FileRouter) Methods() []string {
	return []string{echo.POST}
}

func (fr *FileRouter) Path() string {
	return "/file"
}

func (fr *FileRouter) Handler(c echo.Context) error {
	fmt.Println("handle upload file")
	user := c.FormValue("user")

	file, err := c.FormFile("file")
	if err != nil {
		glog.Error(err)
		return err
	}
	src, err := file.Open()
	if err != nil {
		glog.Error(err)
		return err
	}
	defer src.Close()

	fn := fmt.Sprintf("%s/%s/%d_%s", config.RootDirectory, user, time.Now().Unix(), filepath.Base(file.Filename))
	fmt.Println("upload file ", fn)
	dir := filepath.Dir(fn)
	if exist, _ := PathExists(dir); !exist {
		os.MkdirAll(dir, os.ModePerm)
	}

	dst, err := os.Create(fn)
	if err != nil {
		glog.Error(err)
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		glog.Error(err)
		return err
	}
	return c.String(http.StatusOK, "success")
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

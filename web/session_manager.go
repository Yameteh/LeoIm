package main

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

const (
	SESSION_NAME = "session"
)

type SessionManager struct {
	st *sessions.CookieStore
}

func NewSessionManager() *SessionManager {
	return &SessionManager{sessions.NewCookieStore([]byte(config.SessionKey))}
}

func (sm *SessionManager) GetStore() *sessions.CookieStore {
	return sm.st
}

func (sm *SessionManager) Save(c echo.Context) {
	sess, _ := session.Get(SESSION_NAME, c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   0,
		HttpOnly: true,
	}

	sess.Save(c.Request(), c.Response())
}

func (sm *SessionManager) IsNewSession(c echo.Context) bool {
	var ret bool = false
	sess, _ := session.Get(SESSION_NAME, c)
	if sess != nil {
		ret = sess.IsNew
	}
	return ret
}

func (sm *SessionManager) Remove(c echo.Context) {
	sess, _ := session.Get(SESSION_NAME, c)
	if sess != nil && !sess.IsNew {
		sess.Options.MaxAge = -1
		sess.Save(c.Request(), c.Response())
	}
}

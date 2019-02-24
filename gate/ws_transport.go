package main

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"

	"golang.org/x/net/websocket"
)

type WSTransport struct {
}

func (wst *WSTransport) Listen(domain string, port int) {
	http.Handle("/leoim", websocket.Server{websocket.Config{}, nil, wsHandler})
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", domain, port), nil); err != nil {
		glog.Error(err)
	}
}

func wsHandler(ws *websocket.Conn) {
	ua := uaManager.NewUserAgent(ws)
	ua.Run()
	<-ua.Closer
}

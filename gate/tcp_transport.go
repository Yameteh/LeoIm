package main

import (
	"fmt"
	"github.com/golang/glog"
	"net"
)

type TCPTransport struct {
}

func (ttp *TCPTransport) Listen(domain string, port int) {
	glog.Infof("gate tcp transport listen %s:%d\n", domain, port)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", domain, port))
	if err != nil {
		glog.Error("tcp transport listen error : ", err)
		return
	} else {
		for {
			if conn, err := listener.Accept(); err == nil {
				ua := uaManager.NewUserAgent(conn)
				ua.Run()
			}

		}
	}
}

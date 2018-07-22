package main

import (
	"net"
	"github.com/golang/glog"
	"fmt"
)

type TCPTransport struct {

}



func (ttp *TCPTransport) Listen(domain string, port int) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d",domain,port))
	if err != nil {
		glog.Error("tcp transport listen error : ", err)
		return
	} else {
		for {
			if conn, err := listener.Accept(); err == nil {
				ua := NewUserAgent(conn)
				ua.Run()
			}

		}
	}
}
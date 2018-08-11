package main

import "fmt"

type VideoStream struct {
	Net   string
	Addr  string
	Port  int
	Codec *AudioCodec
	conn  StreamConn
}

func NewVideoStream(net string,addr string,port int) *VideoStream {
	return &VideoStream{net, addr, port, nil,nil}

}

func (as *VideoStream) Record() {
	go func() {
		if as.Net == "tcp" {
			c := new(StreamUdpConn)
			c.Create(as.Addr,as.Port)
			as.conn = c
		}else if as.Net == "udp" {
			c := new(StreamUdpConn)
			c.Create(as.Addr,as.Port)
			as.conn = c
		}

		out := make([]byte,1024)
		for {
			as.conn.Read(out)
			fmt.Println(string(out))
		}
	}();
}

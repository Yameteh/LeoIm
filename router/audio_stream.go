package main

import "fmt"

type AudioStream struct {
	Net   string
	Addr  string
	Port  int
	Codec *AudioCodec
	conn  StreamConn
}

func NewAudioStream(net string, addr string, port int) *AudioStream {
	return &AudioStream{net, addr, port, nil, nil}
}

func (as *AudioStream) Record() {
	go func() {
		var err error
		if as.Net == "tcp" {
			c := new(StreamUdpConn)
			err = c.Create(as.Addr, as.Port)
			as.conn = c
		} else if as.Net == "udp" {
			c := new(StreamUdpConn)
			err = c.Create(as.Addr, as.Port)
			as.conn = c
		}

		if err != nil {
			fmt.Println(err)
		} else {
			out := make([]byte, 1024)
			for {

				as.conn.Read(out)
				fmt.Println(string(out))
			}
		}

	}();
}


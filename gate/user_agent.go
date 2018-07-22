package main

import (
	"net"
	"github.com/satori/go.uuid"
	"fmt"
)

const (
	USER_STATUS_ONLINE = iota
	USER_STATUS_OFFLINE
)

type UserAgent struct {
	Uuid   string
	Conn   net.Conn
	Codec  *ProtocolCodec
	Writer chan *Protocol
	Status uint
}

func NewUserAgent(conn net.Conn) *UserAgent {
	uuid, _ := uuid.NewV4()
	return &UserAgent{
		Uuid:uuid.String(),
		Conn:conn,
		Codec:NewProtocolCodec(conn),
		Writer:make(chan *Protocol),
		Status:USER_STATUS_ONLINE}
}

func (ua *UserAgent) Run() {
	go func() {
		readChan := ua.Codec.Decode()
		for {
			select {
			case r := <-readChan:
				fmt.Println("receive data ",r)
			case w := <-ua.Writer:
				ua.Codec.Encode(w)
			}
		}
	}()
}





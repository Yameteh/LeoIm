package main

import (
	"fmt"
	"encoding/json"
	"github.com/golang/glog"
)

type RouterRpcServer struct {
}

func NewRouterRpcServer() *RouterRpcServer {
	return &RouterRpcServer{}
}

func (rrs *RouterRpcServer) HandleMessage(msg *Message, ret *int) error {
	fmt.Println("handle message ",msg)
	switch msg.Type {
	case 2:
		m := new(MessageBody)
		err := json.Unmarshal([]byte(msg.Body),m)
		if err != nil {
			return err
		}

		





	}
	return nil
}

package main

import "fmt"

type RouterRpcServer struct {
}

func NewRouterRpcServer() *RouterRpcServer {
	return &RouterRpcServer{}
}

func (rrs *RouterRpcServer) HandleMessage(msg *Message, ret *int) error {
	fmt.Println("handle message ",msg)
	
	return nil
}

package main

type RouterRpcServer struct {
}

func NewRouterRpcServer() *RouterRpcServer {
	return &RouterRpcServer{}
}

func (rrs *RouterRpcServer) HandleMessage(msg *Message, ret *int) error {

	return nil
}

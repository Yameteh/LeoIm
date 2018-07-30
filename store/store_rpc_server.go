package main

import "fmt"

type StoreRpcServer struct {
}

func NewStoreRpcServer() *RpcStoreServer {
	return &RpcStoreServer{}
}

func (rs *StoreRpcServer) QueryUser(uuid string, user *User) error {
	fmt.Println("query user uuid ", uuid)
	*user = User{"yaoguoju", "test"}
	return nil
}

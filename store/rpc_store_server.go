package main

import "fmt"

type RpcStoreServer struct {

}

func NewRpcStoreServer() *RpcStoreServer{
	return &RpcStoreServer{}
}

func (rs *RpcStoreServer) QueryUser(uuid string,user *User) error {
	fmt.Println("query user uuid ",uuid)
	*user = User{"yaoguoju","test"}
	return nil
}

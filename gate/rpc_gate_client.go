package main

import (
	"net/rpc"
	"fmt"
	"github.com/golang/glog"
)

type StoreClient struct {
	client *rpc.Client
}

func NewStoreClient() *StoreClient {
	addr := fmt.Sprintf("%s:%d", config.StoreDomain, config.StorePort)
	c, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		glog.Error("dial store server error : ", err)
		return nil
	} else {
		return &StoreClient{c}
	}
}

func (sc *StoreClient) QueryUser(uuid string)  *User {
	var u User
	sc.client.Call("RpcStoreServer.QueryUser", uuid, &u)
	fmt.Println(u.Account)
	return &u
}
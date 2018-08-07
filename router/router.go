package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/rpc"

	"github.com/golang/glog"
)

const (
	ROUTER_CONFIG_INI = "router.ini"
)

var config *Config
var gateManager *GateManager
var storeManager *StoreManager

func main() {
	flag.Parse()
	defer glog.Flush()

	//parse gate config
	config = NewConfig(ROUTER_CONFIG_INI)
	if config != nil {
		config.Parse()
	} else {
		glog.Error("router config ini missed")
		return
	}

	gateManager = NewGateManager()
	storeManager = NewStoreManager()
	err := storeManager.Init()
	if err != nil {
		glog.Error(err)
	}
	setupRouterRpcServer()
	setupWebServer()
	c := make(chan interface{})
	<-c
}

func setupWebServer() {
	webServer := NewRouterWebServer()
	webServer.Serve()
}

func setupRouterRpcServer() {
	go func() {
		rpcSever := NewRouterRpcServer()
		rpc.Register(rpcSever)
		rpc.HandleHTTP()
		addr := fmt.Sprintf("%s:%d", config.Domain, config.Port)
		l, _ := net.Listen("tcp", addr)
		http.Serve(l, nil)
		glog.Info("router rpc server start with address ",addr)
	}()

}

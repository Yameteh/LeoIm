package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/rpc"

	"github.com/golang/glog"
	"time"
)

const (
	ROUTER_CONFIG_INI = "router.ini"
)

var config *Config

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
	rpcSever := NewRouterRpcServer()
	rpc.Register(rpcSever)
	rpc.HandleHTTP()
	addr := fmt.Sprintf("%s:%d", config.Domain, config.Port)
	l, _ := net.Listen("tcp", addr)
	http.Serve(l, nil)

}

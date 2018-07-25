package main

import (
	"flag"
	"github.com/golang/glog"
	"net/rpc"
	"net"
	"fmt"
	"net/http"
)

const (
	STORE_CONFIG_INI = "store.ini"
)

var config *Config

func main() {
	flag.Parse()
	defer glog.Flush()

	config = NewConfig(STORE_CONFIG_INI)

	if config != nil {
		config.Parse()
	} else {
		glog.Error("store config ini missed")
		return
	}

	rpcSever := NewRpcStoreServer()
	rpc.Register(rpcSever)
	rpc.HandleHTTP()
	addr := fmt.Sprintf("%s:%d", config.Domain, config.Port)
	l, _ := net.Listen("tcp", addr)
	http.Serve(l, nil)

}



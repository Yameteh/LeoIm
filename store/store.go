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

	rpcSever := NewStoreRpcServer()
	rpc.Register(rpcSever)
	rpc.HandleHTTP()
	addr := fmt.Sprintf("%s:%d", config.Domain, config.Port)
	l, _ := net.Listen("tcp", addr)
	http.Serve(l, nil)

}

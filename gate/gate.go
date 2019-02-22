package main

import (
	"flag"

	"fmt"
	"net"
	"net/http"
	"net/rpc"

	"github.com/golang/glog"
)

/**
 * gate ini format
 * [listen]
 * domain = localhost
 * port = 8979
 */
const (
	GATE_CONFIG_INI = "gate.ini"
)

var config *Config
var uaManager *AgentManager
var routerManager *RouterManager

func main() {
	//init log
	flag.Parse()
	defer glog.Flush()

	//parse gate config
	config = NewConfig(GATE_CONFIG_INI)
	if config != nil {
		config.Parse()
	} else {
		glog.Error("gate config ini missed")
		return
	}

	uaManager = NewAgentManager()
	routerManager = NewRouterManager(len(config.RouterServer))

	setupGateRpcServer()

	t := new(TCPTransport)
	t.Listen(config.Domain, config.Port)

}

func setupGateRpcServer() {
	go func() {
		rpcSever := NewGateRpcServer()
		rpc.Register(rpcSever)
		rpc.HandleHTTP()
		addr := fmt.Sprintf("%s:%d", config.Domain, config.RpcPort)
		l, err := net.Listen("tcp", addr)
		if err != nil {
			glog.Error(err)
		} else {
			http.Serve(l, nil)
			glog.Info("gate rpc server start with address ", addr)
		}
	}()

}

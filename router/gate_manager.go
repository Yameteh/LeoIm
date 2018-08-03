package main

import (
	"net/rpc"
	"github.com/golang/glog"
)

type GateManager struct {

}

func NewGateManager() *GateManager {
	return &GateManager{}
}

func (gm *GateManager) PublishProtocol(tp *ToProtocol) {
	for _, d := range config.GateServer {
		client, err := rpc.DialHTTP("tcp", d)
		if err != nil {
			glog.Error(err)
			continue
		}
		var ret int
		err = client.Call("GateRpcServer.PublishProtocol", tp, &ret)
		if err != nil {
			glog.Error(err)
			continue
		}
	}
}


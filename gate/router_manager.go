package main

import (
	"net/rpc"
	"github.com/golang/glog"
)

type RouterManager struct {
	RouterCount int
	curIndex    int
	lastIndex   int
}

func NewRouterManager(count int) *RouterManager {
	return &RouterManager{count, 0, -1}
}

func (rm *RouterManager) ChangeCurIndex() {
	rm.curIndex ++
	if rm.curIndex >= rm.RouterCount {
		rm.curIndex = 0
	}

}

func (rm *RouterManager) PublishMessage(p *Protocol) {
	var client *rpc.Client
	var err error
	defer func() {
		if client != nil {
			client.Close()
		}
	}()
	for {
		client, err = rpc.DialHTTP("tcp", config.RouterServer[rm.curIndex])
		if err == nil {
			rm.lastIndex = rm.curIndex
			message := &Message{}
			message.Version = p.Version
			message.Type = p.Type
			message.Body = p.Body
			var ret int
			client.Call("RouterRpcServer.HandleMessage", message, &ret)
			if ret != 0 {
				glog.Error("router server publish response error ")
			}
			rm.ChangeCurIndex()
			return
		} else {
			if rm.curIndex == rm.lastIndex {
				glog.Error("no router server available when publish")
				return
			}
			rm.ChangeCurIndex()
		}

	}
}


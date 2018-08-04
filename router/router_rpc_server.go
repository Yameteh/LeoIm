package main

import (
	"encoding/json"
	"github.com/golang/glog"
	_ "github.com/lib/pq"

	"fmt"
)

type RouterRpcServer struct {
	sm *StoreManager
}

func NewRouterRpcServer() *RouterRpcServer {
	return &RouterRpcServer{NewStoreManager()}
}

func (rrs *RouterRpcServer) SaveMessage(body *MessageBody) {
	err := rrs.sm.Insert(body)
	if err != nil {
		glog.Error(err)
	}
}


func (rrs *RouterRpcServer) HandleMessage(msg *Message, ret *int) error {
	glog.Info("handle message ",msg)
	switch msg.Type {
	case 2:
		m := new(MessageBody)
		err := json.Unmarshal([]byte(msg.Body),m)
		if err != nil {
			return err
		}
		rrs.SaveMessage(m)

		tp := new(ToProtocol)
		tp.To = m.To
		tp.Version = msg.Version
		tp.Type = 4
		tp.Body = fmt.Sprintf("{time:%d}",m.Time)
		tp.Length = uint32(len(tp.Body))
		gateManager.PublishProtocol(tp)





	}
	return nil
}

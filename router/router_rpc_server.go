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
	m := NewStoreManager()
	err := m.Init()
	if err != nil {
		glog.Error(err)
	}
	return &RouterRpcServer{m}
}

func (rrs *RouterRpcServer) saveMessage(body *MessageBody) bool {
	err := rrs.sm.Insert(body)
	if err != nil {
		glog.Error(err)
		return false
	}
	glog.Info("router save message ",body)
	return true
}

func (rrs *RouterRpcServer) HandleMessage(msg *Message, ret *int) error {
	switch msg.Type {
	case 2:
		m := new(MessageBody)
		err := json.Unmarshal([]byte(msg.Body), m)
		if err != nil {
			return err
		}
		if rrs.saveMessage(m) {
			tp := new(ToProtocol)
			tp.To = m.To
			tp.Version = msg.Version
			tp.Type = 4
			tp.Body = fmt.Sprintf("{time:%d}", m.Time)
			tp.Length = uint32(len(tp.Body))
			gateManager.PublishProtocol(tp)
		}else {
			glog.Info("router save message failed ")
			tp := new(ToProtocol)
			tp.To = m.To
			tp.Version = msg.Version
			tp.Type = msg.Type
			tp.Body = msg.Body
			tp.Length = uint32(len(tp.Body))
			glog.Info("router transfer message ",tp)
			gateManager.PublishProtocol(tp)
		}

	}
	return nil
}

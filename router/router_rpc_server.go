package main

import (
	"encoding/json"
	"github.com/golang/glog"
	_ "github.com/lib/pq"

	"fmt"
)

type RouterRpcServer struct {
}

func NewRouterRpcServer() *RouterRpcServer {
	return &RouterRpcServer{}
}

func (rrs *RouterRpcServer) saveMessage(body *MessageBody) bool {
	err := storeManager.Insert(body)
	if err != nil {
		glog.Error(err)
		return false
	}
	glog.Info("router save message ", body)
	return true
}

func (rrs *RouterRpcServer) HandleMessage(msg *Message, ret *int) error {
	go func(msg *Message) {
		switch msg.Type {
		case 2:
			m := new(MessageBody)
			glog.Info("HandleMessage ",msg.Body)
			err := json.Unmarshal(msg.Body, m)
			if err != nil {
				glog.Error(err)
				return;
			}
			if rrs.saveMessage(m) {
				tp := new(ToProtocol)
				tp.To = m.To
				tp.Version = msg.Version
				tp.Type = 4
				sr := &SyncResponse{}
				sr.Time = m.Time
				sr.Server = fmt.Sprintf("%s:%d", config.WebDomain, config.WebPort);
				s, _ := json.Marshal(sr)
				tp.Body = s
				tp.Length = uint32(len(tp.Body))
				gateManager.PublishProtocol(tp)
			} else {
				glog.Info("router save message failed ")
				tp := new(ToProtocol)
				tp.To = m.To
				tp.Version = msg.Version
				tp.Type = msg.Type
				tp.Body = msg.Body
				tp.Length = uint32(len(tp.Body))
				glog.Info("router transfer message ", tp)
				gateManager.PublishProtocol(tp)
			}

		}
	}(msg)
	return nil
}

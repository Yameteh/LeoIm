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

func (rrs *RouterRpcServer) saveMessage(body *MessageBody) {

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
		rrs.saveMessage(m)

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

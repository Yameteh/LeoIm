package main

import (
	"github.com/golang/glog"
	"errors"
)

type GateRpcServer struct {
	
}

func NewGateRpcServer() *GateRpcServer{
	return &GateRpcServer{}
}

func (grs *GateRpcServer) PublishProtocol(tp *ToProtocol,ret *int) error {
	ua := uaManager.getUserAgent(tp.To)
	if ua != nil {
		p := &Protocol{tp.Version,tp.Type,tp.Length,tp.Body}
		glog.Info("publish protocol ",p)
		ua.Writer <- p
		return nil
	}else {
		return errors.New("ua not found ")
	}
}
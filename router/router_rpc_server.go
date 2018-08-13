package main

import (
	"encoding/json"
	"github.com/golang/glog"
	_ "github.com/lib/pq"

	"fmt"
	"github.com/wernerd/GoRTP/src/net/rtp"
	"net"
)

type RouterRpcServer struct {
}

func NewRouterRpcServer() *RouterRpcServer {
	rtp.PayloadFormatMap[96] = &rtp.PayloadFormat{96, rtp.Audio | rtp.Video, 90000, 0, "H264/MB"}
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
				return
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
		case 80:
			sdp := new(StreamSdp)
			err := json.Unmarshal(msg.Body,sdp)
			if err != nil {
				glog.Error(err)
				return
			}


			sdp.InAddr = "172.25.1.137";
			sdp.AudioPort = 10000;
			addr,err := net.ResolveIPAddr("ip",sdp.InAddr)
			fmt.Println(err)
			transportA,err := rtp.NewTransportUDP(addr,sdp.AudioPort,"")
			fmt.Println(err)
			as := NewAudioStream(rtp.NewSession(transportA,transportA))
			as.Record();

			sdp.VideoPort = 10004;
			transportV,_ := rtp.NewTransportUDP(addr,sdp.VideoPort,"")
			vs := NewVideoStream(rtp.NewSession(transportV,transportV))
			vs.Record()

			tp := new(ToProtocol)
			tp.Type = 81
			tp.Version = 1
			tp.To = msg.From
			s , _ := json.Marshal(sdp)
			tp.Body = s
			tp.Length = uint32(len(tp.Body))
			fmt.Println(tp.To)
			gateManager.PublishProtocol(tp)




		}
	}(msg)
	return nil
}

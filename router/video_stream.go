package main

import (
	"github.com/wernerd/GoRTP/src/net/rtp"
)

type VideoStream struct {
	H264p PayloadProcessor
	Session *rtp.Session
}

func NewVideoStream(s *rtp.Session) *VideoStream {
	return &VideoStream{H264p:NewH264Processor(),Session:s}

}


func (vs *VideoStream) Record() {
	go func() {
		rc := vs.Session.CreateDataReceiveChan()

		for p := range rc {
			//p.Print("video")
			vs.H264p.Process(p)
		}
	}()
	vs.Session.StartSession()
}

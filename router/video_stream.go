package main

import (
	"github.com/wernerd/GoRTP/src/net/rtp"
	"github.com/nareix/joy4/cgo/ffmpeg"
)

type VideoStream struct {
	decoder *ffmpeg.VideoDecoder
	Session *rtp.Session
}

func NewVideoStream(s *rtp.Session) *VideoStream {
	return &VideoStream{Session:s}

}

func (vs *VideoStream) initDecoder() {

}

func (vs *VideoStream) Record() {
	go func() {
		rc := vs.Session.CreateDataReceiveChan()

		for p := range rc {
			p.Print("video")
		}
	}()
	vs.Session.StartSession()
}

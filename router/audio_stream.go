package main

import (
	"github.com/wernerd/GoRTP/src/net/rtp"
	"github.com/nareix/joy4/cgo/ffmpeg"
)

type AudioStream struct {
	decoder *ffmpeg.AudioDecoder
	Session *rtp.Session
}

func NewAudioStream(s *rtp.Session) *AudioStream {
	return &AudioStream{Session:s}
}

func(as *AudioStream) initDecoder(){
}


func (as *AudioStream) Record() {
	go func() {
		rc := as.Session.CreateDataReceiveChan()
		for p := range rc {
			p.Print("audio")
		}
	}()
	as.Session.StartSession()
}


package main

import (
	"bufio"
	"github.com/wernerd/GoRTP/src/net/rtp"
	"os"
)

type AudioStream struct {
	Session *rtp.Session
}

func NewAudioStream(s *rtp.Session) *AudioStream {
	return &AudioStream{Session: s}
}

func (as *AudioStream) initDecoder() {
}

func (as *AudioStream) Record() {
	go func() {
		rc := as.Session.CreateDataReceiveChan()
		fileObj, _ := os.OpenFile("test.amr", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		defer fileObj.Close()
		wirteObj := bufio.NewWriter(fileObj)
		header := []byte{0x23, 0x21, 0x41, 0x4d, 0x52, 0x0A}
		wirteObj.Write(header)
		for p := range rc {
			wirteObj.Write(p.Payload())
		}
	}()
	as.Session.StartSession()
}

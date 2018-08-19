package main

import (
	"github.com/wernerd/GoRTP/src/net/rtp"
	"fmt"
	"github.com/nareix/joy4/format/mp4"
	"github.com/nareix/joy4/codec/h264parser"
	"github.com/golang/glog"
	"os"
	"github.com/nareix/joy4/av"
	"time"
)

type VideoStream struct {
	H264p   PayloadProcessor
	Session *rtp.Session
	Muxer   *mp4.Muxer
	sps     []byte
	pps     []byte
	muxing  bool
}

func NewVideoStream(s *rtp.Session) *VideoStream {
	vs := new(VideoStream)
	vs.H264p = NewH264Processor(PayloadReceived(vs.NALReceived))
	vs.Session = s
	vs.muxing = false
	return vs
}

func (vs *VideoStream)NALReceived(n SingleUnit) {
	switch n.NUT() {
	case 7:
		fmt.Println("sps received")
		vs.sps = n.Payload()[1:]
		vs.checkSPSandPPS()
	case 8:
		fmt.Println("pps received")
		vs.pps = n.Payload()[1:]
		vs.checkSPSandPPS()
	case 5:
		if vs.muxing {
			p := av.Packet{}
			p.IsKeyFrame = true
			p.Data = n.Payload()[1:]
			p.Idx = 0
			p.Time =  time.Duration(time.Now().Unix())
			p.CompositionTime =  time.Duration(5*1000)
			err := vs.Muxer.WritePacket(p)
			if err != nil {
				glog.Error(err)
			}
			err = vs.Muxer.WriteTrailer()
			if err != nil {
				glog.Error(err)
			}

		}

	case 1:
		if vs.muxing {
			p := av.Packet{}
			p.IsKeyFrame = false
			p.Data = n.Payload()[1:]
			p.Idx = 0
			p.Time =  time.Duration(time.Now().Unix())
			p.CompositionTime = time.Duration(5*1000)

			err := vs.Muxer.WritePacket(p)
			if err != nil {
				glog.Error(err)
			}

			err = vs.Muxer.WriteTrailer()
			if err != nil {
				glog.Error(err)
			}
		}

	}

}

func (vs *VideoStream) checkSPSandPPS() {
	if vs.pps != nil && vs.sps != nil && !vs.muxing {
		fmt.Println("record mp4 ",vs.sps,",",vs.pps)
		codecData, err := h264parser.NewCodecDataFromSPSAndPPS(vs.sps, vs.pps)
		if err != nil {
			glog.Error(err)
		}
		file, err := os.Create("test1.mp4")
		if err != nil {
			glog.Error(err)
		}
		vs.Muxer = mp4.NewMuxer(file)
		vs.Muxer.WriteHeader([]av.CodecData{codecData})
		vs.muxing = true
	} else {
		fmt.Println("wait record mp4")
	}
}

func (vs *VideoStream) Record() {
	go func() {
		rc := vs.Session.CreateDataReceiveChan()

		for p := range rc {
			vs.H264p.Process(p)
		}
	}()
	vs.Session.StartSession()
}

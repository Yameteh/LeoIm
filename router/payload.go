package main

import (
	"fmt"
	"github.com/wernerd/GoRTP/src/net/rtp"
	"bytes"
	"encoding/binary"
	"github.com/golang/glog"
)

type PayloadProcessor interface {
	Close() error
	Process(*rtp.DataPacket) error
}

type H264Processor struct {
	writable chan SingleUnit
	stop     chan bool
	fua      NALUHandler
	pr       PayloadReceived
}

type PayloadReceived func(nu SingleUnit)

func NewH264Processor(p PayloadReceived) PayloadProcessor {
	writable := make(chan SingleUnit)
	stop := make(chan bool)
	fua := NewFUAHandler()
	handler := &H264Processor{writable: writable, stop: stop, fua: fua, pr:p}
	go handler.outputter(handler.writable, handler.stop)
	return handler
}

func (u *H264Processor) Close() error {
	fmt.Println("Cleaning up...")
	u.stop <- true
	return nil
}

func (u *H264Processor) Process(p *rtp.DataPacket) error {
	n := FromRTP(p)
	switch {
	case n.NUT() <= 23:
		u.writable <- SingleUnit{n}
	case n.NUT() == 28:
		u.fua.Handle(n, u.writable)
	case n.NUT() == 24:
		u.handleStapA(n)
	default:
		fmt.Println("Dropped one type ",n.NUT())
	}
	p.FreePacket()
	return nil
}

func (u *H264Processor) handleStapA(n *NALU) {
	p := n.Payload()
	r := bytes.NewReader(p[1:])
	var naluLen uint16
	var nalu []byte
	fmt.Println("stapA ",n.String())
	fmt.Println("stapA ",p)
	for {
		err := binary.Read(r,binary.BigEndian,&naluLen)
		if err != nil {
			glog.Error(err)
			break
		}
		fmt.Println("nalu len ",naluLen)
		nalu = make([]byte,naluLen)
		err = binary.Read(r,binary.BigEndian,nalu)
		if err != nil {
			glog.Error(err)
			break
		}
		n := FromBytes(nalu,n.Seq(),n.TS())
		u.writable <- SingleUnit{n}
	}

}

func (u *H264Processor) outputter(writable chan SingleUnit, stop chan bool) {
	for {
		select {
		case nalu := <-writable:
			if u.pr != nil {
				u.pr(nalu)
			}
		case <-stop:
			return
		}
	}
}

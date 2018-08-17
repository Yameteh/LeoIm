package main

import (
	"fmt"
	"github.com/wernerd/GoRTP/src/net/rtp"
)

type PayloadProcessor interface {
	Close() error
	Process(*rtp.DataPacket) error
}

type H264Processor struct {
	writable chan SingleUnit
	stop chan bool
	fua NALUHandler
}

func NewH264Processor() PayloadProcessor {
	writable := make(chan SingleUnit)
	stop := make(chan bool)
	fua := NewFUAHandler()
	handler := &H264Processor{ writable: writable, stop: stop, fua: fua}
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
	fmt.Println("nalu type ",n.NUT())
	switch {
	case n.NUT() <= 23:
		u.writable <- SingleUnit{n}
	case n.NUT() == 28:
		u.fua.Handle(n, u.writable)
	default:
		fmt.Println("Dropped one")
	}
	p.FreePacket()
	return nil
}

func (u *H264Processor) outputter(writable chan SingleUnit, stop chan bool) {
	for {
		select {
		case nalu := <-writable:
			//u.writer.Write([]byte{0x00, 0x00, 0x00, 0x01})
			//_, e := u.writer.Write(nalu.Payload())
			//if e != nil {
			//	fmt.Println("Write error")
			//}

			fmt.Println(nalu)
		case <-stop:
			return
		}
	}
}

package main

import (
	"fmt"
)

type NALUHandler interface {
	NALUTypes() []int
	Handle(nalu *NALU, writable chan SingleUnit) error
}

type FUAHandler struct {
	naluTypes []int
	buffer    []FUAUnit
}

func (f *FUAHandler) NALUTypes() []int {
	return f.naluTypes
}

func (f *FUAHandler) Handle(nalu *NALU, output chan SingleUnit) error {
	fua := FUAUnit{nalu}

	if len(f.buffer) == 0 && fua.Start() {
		f.buffer = append(f.buffer, fua)
		return nil
	}

	if len(f.buffer) > 0 && !fua.Start() && nalu.Seq() - f.buffer[len(f.buffer) - 1].Seq() == 1 {
		f.buffer = append(f.buffer, fua)

		if fua.End() {
			output <- toSingleUnit(f.buffer)
			f.buffer = make([]FUAUnit, 0)
		}
		return nil
	}

	fmt.Println("Bad FU-A Packet")
	f.buffer = make([]FUAUnit, 0)
	return nil
}

func NewFUAHandler() NALUHandler {
	return &FUAHandler{naluTypes: []int{28}, buffer: make([]FUAUnit, 0)}
}

func toSingleUnit(fuas []FUAUnit) SingleUnit {
	first := 0xE0 & fuas[0].Payload()[0]
	head := int8(first) + fuas[0].PayNUT()
	suPayload := make([]byte, 0)
	suPayload = append(suPayload, byte(head))

	for _, value := range fuas {
		suPayload = append(suPayload, value.FPayload()...)
	}
	su := SingleUnit{FromBytes(suPayload, fuas[0].Seq(), fuas[0].TS())}
	return su
}

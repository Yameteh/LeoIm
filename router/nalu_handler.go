package main

import (
	"fmt"
)

type NALUHandler interface {
	NALUTypes() []int //NALU types the handler will accept
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

	// If the buffer this empty and this is a start packet...
	if len(f.buffer) == 0 && fua.Start() {
		f.buffer = append(f.buffer, fua)
		return nil
	}

	// If the buffer is not empty and this is not a start packet and the
	// sequence number is consecutive
	if len(f.buffer) > 0 && !fua.Start() && nalu.Seq()-f.buffer[len(f.buffer)-1].Seq() == 1 {
		f.buffer = append(f.buffer, fua)

		// If this was also an End packet, construct a SingleUnit NAL and write
		if fua.End() {
      output <- toSingleUnit(f.buffer)
      f.buffer = make([]FUAUnit, 0)
		}
		return nil
	}

	// If we made it here without returning, this was an invalid packet
	fmt.Println("Bad FU-A Packet")
  f.buffer = make([]FUAUnit, 0)
	return nil
}

func NewFUAHandler() NALUHandler {
	return &FUAHandler{naluTypes: []int{28}, buffer: make([]FUAUnit, 0)}
}

func toSingleUnit(fuas []FUAUnit) SingleUnit {
  // Get the first 3 bits of the NALU header
  first := 0xE0 & fuas[0].Payload()[0]
  // Combined with FU Payload NUT gives the new NALU header
  head := int8(first) + fuas[0].PayNUT()
  suPayload := make([]byte, 0)
  suPayload = append(suPayload, byte(head))

  for _, value := range fuas {
    suPayload = append(suPayload, value.FPayload()...)
  }
  su := SingleUnit{FromBytes(suPayload, fuas[0].Seq(), fuas[0].TS())}
  return su
}

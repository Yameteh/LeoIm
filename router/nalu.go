package main

import (
	"fmt"
	"github.com/wernerd/GoRTP/src/net/rtp"
)

type NALU struct {
	forbidden bool   // forbidden_zero_bit
	nri       int8   // nal_ref_idc
	nut       int8   // nal_unit_type
	seq       uint16 // sequence_number (RTP)
	ts        uint32 // timestamp (RTP)
	payload   []byte // entire payload from RTP packet
}

const (
	headIdx = 0
	fMask = 0x80
	nriMask = 0x60
	typeMask = 0x1f
)

func FromRTP(rtp *rtp.DataPacket) *NALU {
	return FromBytes(rtp.Payload(), rtp.Sequence(), rtp.Timestamp())
}

func FromBytes(payload []byte, seq uint16, ts uint32) *NALU {
	return &NALU{payload: payload, seq: seq, ts: ts}
}

func (n *NALU) Payload() []byte {
	return n.payload
}

func (n *NALU) Forbidden() bool {
	return n.Payload()[headIdx] & fMask >> 7 != 0
}

func (n *NALU) NRI() int8 {
	return int8(n.Payload()[headIdx] & nriMask >> 5)
}

func (n *NALU) NUT() int8 {
	return int8(n.Payload()[headIdx] & typeMask)
}

func (n *NALU) Seq() uint16 {
	return n.seq
}

func (n *NALU) TS() uint32 {
	return n.ts
}

func (n *NALU) String() string {
	return fmt.Sprintf("Forbidden: %v, NRI: %b, Type: %v, Seq: %v, TS: %v, Len: %v", n.Forbidden(), n.NRI(), n.NUT(), n.Seq(), n.TS(), len(n.Payload()))
}

type SingleUnit struct {
	*NALU
}

type FUAUnit struct {
	*NALU
}

const (
	fuaHIdx = 1
	fuaHStartMask = 0x80
	fuaHEndMask = 0x40
	fuaHResMask = 0x20
	fuaHNUTMask = 0x1f
	fuaPayIdx = 2
)

func (n *NALU) Start() bool {
	return n.Payload()[fuaHIdx] & fuaHStartMask >> 7 != 0
}

func (n *NALU) End() bool {
	return n.Payload()[fuaHIdx] & fuaHEndMask >> 6 != 0
}

func (n *NALU) Reserved() bool {
	return n.Payload()[fuaHIdx] & fuaHResMask >> 5 != 0
}

func (n *NALU) PayNUT() int8 {
	return int8(n.Payload()[fuaHIdx] & fuaHNUTMask)
}

func (n *NALU) FPayload() []byte {
	return n.Payload()[2:]
}

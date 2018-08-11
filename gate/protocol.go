package main

import (
	"encoding/binary"
	"io"
	"net"

	"github.com/golang/glog"
)

const (
	PROTOCOL_TYPE_AUTH = 0
	PROTOCOL_TYPE_AUTHACK = 1
	PROTOCOL_TYPE_MSG = 2
	PROTOCOL_TYPE_MSGACK = 3
	PROTOCOL_TYPE_MSGSYNC = 4
	PROTOCOL_TYPE_STREAM_RECORD = 80
)

type Protocol struct {
	Version uint8
	Type    uint8
	Length  uint32
	Body    []byte
}

type ProtocolCodec struct {
	Writer io.Writer
	Reader io.Reader
}

func NewProtocolCodec(conn net.Conn) *ProtocolCodec {
	return &ProtocolCodec{conn, conn}
}

func (pc *ProtocolCodec) Decode() (*Protocol, error) {
	protocol := &Protocol{}
	if err := binary.Read(pc.Reader, binary.BigEndian, &protocol.Version); err != nil {
		return protocol, err
	}
	if err := binary.Read(pc.Reader, binary.BigEndian, &protocol.Type); err != nil {
		return protocol, err
	}
	if err := binary.Read(pc.Reader, binary.BigEndian, &protocol.Length); err != nil {
		return protocol, err
	}
	body := make([]byte, protocol.Length)
	if err := binary.Read(pc.Reader, binary.BigEndian, body); err != nil {
		return protocol, err
	}
	protocol.Body = body
	glog.Info("decode protocol ", protocol)
	return protocol, nil
}

func (pc *ProtocolCodec) Encode(p *Protocol) error {
	if err := binary.Write(pc.Writer, binary.BigEndian, p.Version); err != nil {
		return err
	}
	if err := binary.Write(pc.Writer, binary.BigEndian, p.Type); err != nil {
		return err
	}
	if err := binary.Write(pc.Writer, binary.BigEndian, p.Length); err != nil {
		return err
	}
	if err := binary.Write(pc.Writer, binary.BigEndian, p.Body); err != nil {
		return err
	}
	glog.Info("encode protocol ", p)
	return nil
}

package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/golang/glog"
	"net"
	"fmt"
)

const (
	PROTOCOL_TYPE_AUTH = 0
	PROTOCOL_TYPE_AUTHACK = 1
)

type Protocol struct {
	Version uint8
	Type    uint8
	Length  uint32
	Body    string
}

type ProtocolCodec struct {
	Writer *bufio.Writer
	Reader *bufio.Reader
}

func NewProtocolCodec(conn net.Conn) *ProtocolCodec {
	return &ProtocolCodec{bufio.NewWriter(conn), bufio.NewReader(conn)}
}

func (pc *ProtocolCodec) Decode() chan *Protocol {
	out := make(chan *Protocol)
	go func() {
		for {
			p, err := pc.Reader.Peek(6)
			if err != nil {
				glog.Warning("protocol decode peek error : ", err)
				out <- nil
				return
			}
			protocol := &Protocol{}
			v, t, l := pc.GetHeader(p)
			glog.Infof("decode header [version:%d,type=%d,length=%d]\n", v, t, l)
			protocol.Version = v
			protocol.Type = t
			protocol.Length = l
			pc.Reader.Discard(6)
			b := make([]byte, l)
			buf := bytes.NewBuffer(b)
			buf.Reset()
			for {
				var outSize uint32 = l
				s := pc.Reader.Buffered()
				if uint32(s) >= outSize {
					if o, err := pc.Reader.Peek(int(outSize)); err == nil {
						binary.Write(buf, binary.BigEndian, o)
						pc.Reader.Discard(int(outSize))
					} else {
						glog.Warning(err)
						out <- nil
						return
					}
					break
				} else {
					if o, err := pc.Reader.Peek(s); err == nil {
						binary.Write(buf, binary.BigEndian, o)
						outSize = outSize - uint32(s)
						pc.Reader.Discard(int(s))
					} else {
						glog.Warning(err)
						out <- nil
						return
					}
				}

			}
			protocol.Body = buf.String()
			glog.Infof("decode body %s\n", protocol.Body)
			out <- protocol

		}
	}()
	return out
}

func (pc *ProtocolCodec) Encode(p *Protocol) {
	out := pc.ToBytes(p)
	len,err := pc.Writer.Write(out)
	fmt.Println("write len ",len)
	if err != nil {
		fmt.Println("write error ",err)
	}
}

func (pc *ProtocolCodec) GetHeader(b []byte) (uint8, uint8, uint32) {
	var version, types uint8
	var length uint32
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.BigEndian, &version)
	binary.Read(buf, binary.BigEndian, &types)
	binary.Read(buf, binary.BigEndian, &length)
	return version, types, length
}

func (pc *ProtocolCodec) ToBytes(p *Protocol) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, p.Version)
	binary.Write(buf, binary.BigEndian, p.Type)
	binary.Write(buf, binary.BigEndian, p.Length)
	binary.Write(buf, binary.BigEndian, []byte(p.Body))
	return buf.Bytes()
}
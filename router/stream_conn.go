package main

import (
	"fmt"
	"net"
)

type StreamConn interface {
	Create(addr string, port int) error
	Read(out []byte) (int, error)
}

type StreamUdpConn struct {
	conn *net.UDPConn
}

func (suc *StreamUdpConn) Create(addr string, port int) error {
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", addr, port))
	if err == nil {
		c, err := net.ListenUDP("udp", udpAddr)
		if err == nil {
			suc.conn = c
		}
	}
	fmt.Println(err)
	return err
}

func (suc *StreamUdpConn) Read(out []byte) (int, error) {
	len, _, err := suc.conn.ReadFromUDP(out)
	return len, err
}

type StreamTcpConn struct {
	conn net.Conn
}

func (stc *StreamTcpConn) Create(addr string, port int) error {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err == nil {
		c, err := l.Accept()
		if err == nil {
			stc.conn = c
		}
	}
	return err
}

func (stc *StreamTcpConn) Read(out []byte) (int, error) {
	l, err := stc.conn.Read(out)
	return l, err
}

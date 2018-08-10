package main

type Message struct {
	Version uint8
	Type    uint8
	Body    []byte
}

type ToProtocol struct {
	To string
	Version uint8
	Type uint8
	Length uint32
	Body []byte
}
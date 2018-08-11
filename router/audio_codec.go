package main

type AudioCodec interface {
	encode(in []byte) []byte
	decode(out []byte) []byte
}

package main

import (
	"github.com/nareix/joy4/format/mp4"
)

type LeoStream struct {
	muxer  *mp4.Muxer
	parser PayloadProcessor
}

func NewLeoStream() {

}

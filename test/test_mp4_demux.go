package main

import (
	"github.com/nareix/joy4/format/mp4"
	"os"
	"fmt"
	"github.com/nareix/joy4/av"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Command args error")
		return
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	dm := mp4.NewDemuxer(file)

	streams,err := dm.Streams()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, stream := range streams {
		if stream.Type().IsAudio() {
			astream := stream.(av.AudioCodecData)
			fmt.Println(astream.Type(), astream.SampleRate(), astream.SampleFormat(), astream.ChannelLayout())
		} else if stream.Type().IsVideo() {
			vstream := stream.(av.VideoCodecData)
			fmt.Println(vstream.Type(), vstream.Width(), vstream.Height())
		}
	}

	for {
		var pkt av.Packet
		var err error
		if pkt, err = dm.ReadPacket(); err != nil {
			break
		}
		if pkt.IsKeyFrame {
			fmt.Println("pkt", streams[pkt.Idx].Type(), "len", len(pkt.Data), "keyframe", pkt.IsKeyFrame,
				"CompositionTime ",pkt.CompositionTime,"Time",pkt.Time)
		}

	}
}


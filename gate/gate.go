package main

import (
	"flag"
	"github.com/golang/glog"
	"fmt"
)

/**
 * gate ini format
 * [listen]
 * domain = localhost
 * port = 8979
 */
const (
	GATE_CONFIG_INI = "gate.ini"
)

var config *Config

func main() {
	//init log
	flag.Parse()
	defer glog.Flush()

	//parse gate config
	config = NewConfig(GATE_CONFIG_INI)
	if config != nil {
		config.Parse()
	}else {
		glog.Error("gate config ini missed")
		return
	}

	t := new(TCPTransport)
	t.Listen(config.Domain,config.Port)
}

func testProtocol() {
	p := &Protocol{}
	p.Version = 1
	p.Type = 2
	p.Length = 2123
	p.Body = "abc"

	codec := &ProtocolCodec{}
	bytes := codec.ToBytes(p)
	fmt.Println(bytes)
	v,t,l := codec.GetHeader(bytes)
	fmt.Printf("version %d type %d length %d",v,t,l)
}



package main

import (
	"flag"

	"github.com/golang/glog"
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
var uaManager *AgentManager

func main() {
	//init log
	flag.Parse()
	defer glog.Flush()

	//parse gate config
	config = NewConfig(GATE_CONFIG_INI)
	if config != nil {
		config.Parse()
	} else {
		glog.Error("gate config ini missed")
		return
	}

	uaManager = NewAgentManager()

	t := new(TCPTransport)
	t.Listen(config.Domain, config.Port)

}

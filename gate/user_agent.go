package main

import (
	"net"
	"github.com/satori/go.uuid"
	"sync"
	"github.com/golang/glog"
	"fmt"
)

const (
	USER_STATUS_ONLINE = iota
	USER_STATUS_OFFLINE
)

type AgentManager struct {
	sync.RWMutex
	agents map[string]*UserAgent
}

func NewAgentManager() *AgentManager {
	return &AgentManager{agents:make(map[string]*UserAgent)}
}

func (am *AgentManager) NewUserAgent(conn net.Conn) *UserAgent {
	uuid, _ := uuid.NewV4()
	ua := &UserAgent{
		Uuid:uuid.String(),
		Conn:conn,
		Codec:NewProtocolCodec(conn),
		Writer:make(chan *Protocol),
		Status:USER_STATUS_ONLINE}
	am.putUserAgent(ua)
	return ua
}

func (am *AgentManager) putUserAgent(ua *UserAgent) {
	am.Lock()
	defer am.Unlock()
	glog.Infof("user agent %s added",ua.Uuid)
	am.agents[ua.Uuid] = ua
}

func (am *AgentManager) getUserAgent(uuid string) *UserAgent {
	am.RLock()
	defer am.RUnlock()
	ua, _ := am.agents[uuid]
	return ua
}

func (am *AgentManager) delUserAgent(ua *UserAgent)  {
	am.Lock()
	defer am.Unlock()
	glog.Infof("user agent %s deleted",ua.Uuid)
	delete(am.agents,ua.Uuid)
}

type UserAgent struct {
	Uuid   string
	Conn   net.Conn
	Codec  *ProtocolCodec
	Writer chan *Protocol
	Status uint
}

func (ua *UserAgent) Run() {
	go func() {
		readChan := ua.Codec.Decode()
		for {
			select {
			case r := <-readChan:
				fmt.Println("receive data ",r)
				u := storeClient.QueryUser("123")
				fmt.Println(u)
				if r == nil {
					ua.Conn.Close()
					uaManager.delUserAgent(ua)
				}
			case w := <-ua.Writer:
				ua.Codec.Encode(w)
			}
		}
	}()
}





package main

import (
	"net"
	"github.com/satori/go.uuid"
	"sync"
	"github.com/golang/glog"
	"fmt"
	"time"
	"encoding/json"
)

const (
	USER_STATUS_UNKOWN = iota
	USER_STATUS_ONLINE
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
		Writer:make(chan *Protocol,2),
		Status:USER_STATUS_UNKOWN}
	am.putUserAgent(ua)
	return ua
}

func (am *AgentManager) putUserAgent(ua *UserAgent) {
	am.Lock()
	defer am.Unlock()
	glog.Infof("user agent %s added", ua.Uuid)
	am.agents[ua.Uuid] = ua
}

func (am *AgentManager) getUserAgent(uuid string) *UserAgent {
	am.RLock()
	defer am.RUnlock()
	ua, _ := am.agents[uuid]
	return ua
}

func (am *AgentManager) delUserAgent(ua *UserAgent) {
	am.Lock()
	defer am.Unlock()
	glog.Infof("user agent %s deleted", ua.Uuid)
	delete(am.agents, ua.Uuid)
}

type UserAgent struct {
	Uuid   string
	Conn   net.Conn
	Codec  *ProtocolCodec
	Writer chan *Protocol
	Status uint
	Nonce  string
	User   *User
}

func (ua *UserAgent) HandleProtocol(p *Protocol) {
	if p.Type == PROTOCOL_TYPE_AUTH {
		var req AuthRequest
		err := json.Unmarshal([]byte(p.Body), &req)
		if err != nil {
			fmt.Println("auth req parse error :", err)
			return
		}
		ua.Auth(req)
	} else {

	}
}

func (ua *UserAgent) Auth(r AuthRequest) {
	glog.Infof("auth request [user:%s response:%s]\n", r.User, r.Response)
	if r.User != "" && r.Response == "" {
		rsp := new(AuthResponse)
		rsp.Method = config.AuthMethod
		n, _ := uuid.NewV1()
		rsp.Code = 401
		rsp.Nonce = n.String()
		ua.Nonce = rsp.Nonce
		p := new(Protocol)
		p.Version = 1
		p.Type = PROTOCOL_TYPE_AUTHACK
		j, _ := json.Marshal(rsp)
		p.Body = string(j)
		p.Length = uint32(len(p.Body))
		ua.Writer <- p
	} else if r.User != "" && r.Response != "" {
		in := r.User + ":" + ua.Nonce + ":" + "1234"
		if r.Response == GetMd5(in) {
			rsp := new(AuthResponse)
			rsp.Code = 200
			p := new(Protocol)
			p.Version = 1
			p.Type = PROTOCOL_TYPE_AUTHACK
			j,_:= json.Marshal(rsp)
			p.Body = string(j)
			p.Length = uint32(len(p.Body))
			ua.Writer <- p
		} else {
			rsp := new(AuthResponse)
			rsp.Code = 402
			p := new(Protocol)
			p.Version = 1
			p.Type = PROTOCOL_TYPE_AUTHACK
			j,_:= json.Marshal(rsp)
			p.Body = string(j)
			p.Length = uint32(len(p.Body))
			ua.Writer <- p
		}
	}
}


func (ua *UserAgent) Run() {
	go func() {
		readChan := ua.Codec.Decode()
		for {
			select {
			case r := <-readChan:
				fmt.Println("receive data ", r)
				if r == nil {
					ua.Conn.Close()
					uaManager.delUserAgent(ua)
					return
				}else {
					ua.HandleProtocol(r)

				}
			case w := <-ua.Writer:
				fmt.Println("write data ",w)
				ua.Codec.Encode(w)

			case <-time.After(20 * time.Second):
				if ua.Status == USER_STATUS_UNKOWN {
					fmt.Println("auth timeout")
					ua.Conn.Close()
					uaManager.delUserAgent(ua)
					return
				}

			}
		}
	}()
}





package main

import (
	"encoding/json"
	"github.com/golang/glog"
	"github.com/satori/go.uuid"
	"net"
	"sync"
	"time"
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
	return &AgentManager{agents: make(map[string]*UserAgent)}
}

func (am *AgentManager) NewUserAgent(conn net.Conn) *UserAgent {
	ua := &UserAgent{
		Conn:   conn,
		Codec:  NewProtocolCodec(conn),
		Writer: make(chan *Protocol),
		Status: USER_STATUS_UNKOWN}
	return ua
}

func (am *AgentManager) putUserAgent(ua *UserAgent) {
	am.Lock()
	defer am.Unlock()
	glog.Infof("user agent %s added", ua.User.Uuid)
	am.agents[ua.User.Uuid] = ua
}

func (am *AgentManager) getUserAgent(uuid string) *UserAgent {
	am.RLock()
	defer am.RUnlock()
	ua, _ := am.agents[uuid]
	return ua
}

func (am *AgentManager) delUserAgent(ua *UserAgent) {
	if ua.User != nil {
		am.Lock()
		defer am.Unlock()
		glog.Infof("user agent %s deleted", ua.User.Uuid)
		delete(am.agents, ua.User.Uuid)
	}
}

type UserAgent struct {
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
		err := json.Unmarshal(p.Body, &req)
		if err != nil {
			return
		}
		ua.Auth(req)
	} else {
		routerManager.PublishMessage(ua.User.Uuid, p)
	}
}

func (ua *UserAgent) Auth(r AuthRequest) {
	glog.Infof("auth request [user:%s response:%s]\n", r.User, r.Response)
	if r.User != "" && r.Response == "" {
		rsp := new(AuthResponse)
		rsp.Code = AUTHACL_CODE_REAUTH
		rsp.Method = config.AuthMethod
		n, _ := uuid.NewV1()
		rsp.Nonce = n.String()
		ua.Nonce = rsp.Nonce
		p := CreateProtocolMsg(1, PROTOCOL_TYPE_AUTHACK, rsp)
		ua.Writer <- p
	} else if r.User != "" && r.Response != "" {
		user := RedisQueryUser(r.User)
		var in string
		if user != nil {
			in = r.User + ":" + ua.Nonce + ":" + user.Password
		} else {
			in = ""
		}
		if r.Response == in {
			token, _ := uuid.NewV1()
			user.Token = token.String()
			ua.User = user
			ua.Status = USER_STATUS_ONLINE
			uaManager.putUserAgent(ua)
			err := RedisUpdateUser(user)
			if err != nil {
				glog.Error("redis update user error ", err)
			} else {
				rsp := new(AuthResponse)
				rsp.Code = AUTHACL_CODE_OK
				rsp.Token = user.Token
				p := CreateProtocolMsg(1, PROTOCOL_TYPE_AUTHACK, rsp)
				ua.Writer <- p
			}

		} else {
			rsp := new(AuthResponse)
			rsp.Code = AUTHACL_CODE_ERROR
			p := CreateProtocolMsg(1, PROTOCOL_TYPE_AUTHACK, rsp)
			ua.Writer <- p
		}
	}
}

func (ua *UserAgent) Close() {
	ua.Writer <- nil
	ua.Conn.Close()
	ua.Status = USER_STATUS_UNKOWN
	uaManager.delUserAgent(ua)
}

func (ua *UserAgent) Run() {
	go func() {
		for {
			p, err := ua.Codec.Decode()
			if err != nil {
				ua.Close()
				return
			} else {
				ua.HandleProtocol(p)
			}

		}
	}()

	go func() {
		for {
			select {
			case w := <-ua.Writer:
				if w != nil {
					ua.Codec.Encode(w)
				} else {
					glog.Error("ua close write goroutine ")
					return
				}
			}
		}
	}()

	go func() {
		<-time.After(20 * time.Second)
		if ua.Status == USER_STATUS_UNKOWN {
			ua.Close()
		}
	}()

}

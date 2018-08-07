package main

import (
	"fmt"
	"net"
	"bytes"
	"encoding/binary"
	"time"
	"math/rand"

	"bufio"
	"os"
	"strings"
	"encoding/json"
	"net/http"
	"strconv"
	"io/ioutil"
)

type AuthRequest struct {
	User     string
	Response string
}

type AuthResponse struct {
	Code   int
	Nonce  string
	Method string
	Token string
}

func ToBytes(p *Protocol) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, p.Version)
	binary.Write(buf, binary.BigEndian, p.Type)
	binary.Write(buf, binary.BigEndian, p.Length)
	binary.Write(buf, binary.BigEndian, []byte(p.Body))
	return buf.Bytes()
}

type MessageBody struct {
	MsgType int
	From string
	To string
	Time int64
	MimeType string
	Content string
}

type SyncResponse struct {
	Server string
	Time  int64
}


func GetRandomString(length int64) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var i int64
	for i = 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func main() {
	fmt.Println("gate test start")
	conn, err := net.Dial("tcp", "localhost:8979")
	if err != nil {
		fmt.Println("dial error : ", err)
		return
	}
	codec = &ProtocolCodec{conn, conn}
	go func() {
		for {
			p, err := codec.Decode()
			if (err != nil) {
				return
			}
			switch p.Type {
			case 1:
				reAuth(p)
			case PROTOCOL_TYPE_MSGSYNC:
				fmt.Println("msg sync ",p)
				sync(p)
			case PROTOCOL_TYPE_MSG:
				fmt.Println("received msg ",p)
			}
		}
	}()
	r := bufio.NewReader(os.Stdin)
	if r != nil {
		for {
			l, _, err := r.ReadLine()
			if err != nil {
				return
			} else {
				line := string(l)
				cmds := strings.Split(line, " ")
				switch cmds[0] {
				case "auth":
					if len(cmds) == 3 {
						Auth(cmds[1], cmds[2])
					} else {
						fmt.Println("ps : auth xx xx")
					}
				case "msg":
					if len(cmds) == 3 {
						Msg(cmds[1],cmds[2])
					}else {
						fmt.Println("ps : msg xx xxx")
					}
				}

			}
		}
	}
}

func sync(p *Protocol) {
	a := &SyncResponse{}
	json.Unmarshal([]byte(p.Body),a)
	url :=  fmt.Sprintf("http://%s/sync/",a.Server)
	r,err  := http.NewRequest("GET",url,nil)
	if err != nil {
		fmt.Println(err)
	}
	r.Header.Add("Authorization","Basic "+GetBase64(account+":"+token))
	r.Header.Add("SyncTime",strconv.FormatInt(a.Time,10))
	rsp,err  := http.DefaultClient.Do(r)
	if err != nil {
		fmt.Println(err)
	}else {
		fmt.Println(rsp.StatusCode,rsp.ContentLength)
		defer rsp.Body.Close()
		a,err := ioutil.ReadAll(rsp.Body)
		if err == nil {
			fmt.Println(string(a))
		}else {
			fmt.Println(err)
		}

	}
}


var account string
var password string
var token string
var codec *ProtocolCodec

func Auth(user string, pwd string) {
	account = user
	password = pwd
	p := &Protocol{}
	p.Version = 1
	p.Type = 0

	a := &AuthRequest{}
	a.User = user
	a.Response = ""
	s, _ := json.Marshal(a)
	p.Body = string(s)
	p.Length = uint32(len(p.Body))
	err := codec.Encode(p)
	if err != nil {
		fmt.Println(err)
	}
}

func reAuth(in *Protocol) {
	ar := &AuthResponse{}
	json.Unmarshal([]byte(in.Body), ar)
	if ar.Code == 401 {
		fmt.Println("re auth")
		a := &AuthRequest{}
		a.User = account
		r := a.User + ":" + ar.Nonce + ":" + password
		a.Response = r
		fmt.Println("response ",a.Response)
		p := CreateProtocolMsg(1, 0, a)
		err := codec.Encode(p)
		if err != nil {
			fmt.Println(err)
		}
	}else if ar.Code == 200 {
		fmt.Println("auth success token ",ar.Token)
		token = ar.Token
	}else if ar.Code == 402 {
		fmt.Println("auth failed")
	}
}

func Msg(uuid string,content string) {
	message := new(MessageBody)
	message.MsgType = 1
	message.From = account
	message.To = uuid
	message.MimeType = "text/plain"
	message.Time = time.Now().Unix()
	fmt.Println("send time ",message.Time)
	message.Content = content
	p := CreateProtocolMsg(1,PROTOCOL_TYPE_MSG,message)
	err := codec.Encode(p)
	if err != nil {
		fmt.Println(err)
	}
}

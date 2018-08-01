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
)

type AuthRequest struct {
	User     string
	Response string
}

type AuthResponse struct {
	Code   int
	Nonce  string
	Method string
}

func ToBytes(p *Protocol) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, p.Version)
	binary.Write(buf, binary.BigEndian, p.Type)
	binary.Write(buf, binary.BigEndian, p.Length)
	binary.Write(buf, binary.BigEndian, []byte(p.Body))
	return buf.Bytes()
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

				}
			}
		}
	}

	//
	//p := &Protocol{}
	//p.Version = 1
	//p.Type = 3
	//
	//a := &AuthRequest{}
	//a.User="123"
	//a.Response = ""
	//s ,_ := json.Marshal(a)
	//p.Body = string(s)
	//p.Length = uint32(len(p.Body))
	//err = codec.Encode(p)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//c := make(chan interface{})
	//<-c


}

var account string
var password string
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
		fmt.Println("auth success")
	}else if ar.Code == 402 {
		fmt.Println("auth failed")
	}


}


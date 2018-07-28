package main

import (
	"fmt"
	"net"
	"bytes"
	"encoding/binary"
	"time"
	"math/rand"
	"encoding/json"
)

type AuthRequest struct {
	User string
	Response string
}

type AuthResponse struct {
	Nonce string
	Method string
}

func  ToBytes(p *Protocol) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, p.Version)
	binary.Write(buf, binary.BigEndian, p.Type)
	binary.Write(buf, binary.BigEndian, p.Length)
	binary.Write(buf, binary.BigEndian, []byte(p.Body))
	return buf.Bytes()
}

func GetRandomString(length int64) string{
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var i int64
	for i  = 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func main() {
	fmt.Println("gate test start")
	conn, err := net.Dial("tcp", "localhost:8979")
	if err != nil {
		fmt.Println("dial error : ", err)
	}

	//for  i := 1 ; i < 10 ; i++ {
	//	s := GetRandomString(int64(i))
	//	p  := new(Protocol)
	//	p.Version = 1
	//	p.Type = 12
	//	p.Body = s
	//	p.Length = uint32(len(s))
	//	conn.Write(ToBytes(p))
	//}


	codec := &ProtocolCodec{conn,conn}
	go func() {
		for {
			p, err := codec.Decode()
			if (err != nil) {
				return ;
			}
			fmt.Println(p)
		}
	}()

	p := &Protocol{}
	p.Version = 1
	p.Type = 0

	a := &AuthRequest{}
	a.User="123"
	a.Response = ""
	s ,_ := json.Marshal(a)
	p.Body = string(s)
	p.Length = uint32(len(p.Body))
	err = codec.Encode(p)
	if err != nil {
		fmt.Println(err)
	}
	c := make(chan interface{})
	<-c
	

}


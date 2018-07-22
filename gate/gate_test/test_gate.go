package main

import (
	"fmt"
	"net"
	"bytes"
	"encoding/binary"
	"time"
	"math/rand"
)

type Protocol struct {
	Version uint8
	Type    uint8
	Length  uint32
	Body    string
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

	for i := 0; i < 10; i++ {
		p := new(Protocol)
		p.Version = 2
		p.Type = uint8(i)
		s := GetRandomString(int64(i)+10)
		p.Length = uint32(len(s))
		p.Body = s

		fmt.Println(p)
		conn.Write(ToBytes(p))
	}

	c := make(chan int)
	<-c
}


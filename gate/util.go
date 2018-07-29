package main

import "crypto/md5"

func GetMd5(in string) string {
 	h := md5.New()
	h.Write([]byte(in))
	return string(h.Sum(nil))
}
package main

import (
	"log"
	"net"
)

func main() {
	//addr, err := net.ResolveTCPAddr("tcp", "localhost:9003")
	addr, err := net.ResolveTCPAddr("tcp", "www.baidu.com:8080")
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Print(addr)
}

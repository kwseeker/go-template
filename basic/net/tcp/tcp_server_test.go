package tcp

import (
	"fmt"
	"net"
	"testing"
)

// TCP服务器测试
func TestTCP(t *testing.T) {
	fmt.Println("server start...")
	listen, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	for {
		fmt.Println("loop ...")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go process(conn)
	}
}

func process(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 512)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read err:", err)
			return
		}
		fmt.Println("read:", string(buf[:n]))
	}
}

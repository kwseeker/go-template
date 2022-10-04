package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
)

func main() {
	url := "ws://localhost:9002/echo"
	origin := "http://localhost/"

	conn, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := conn.Write([]byte("{\"userId\":10001,\"code\":1001,\"data\":{}}")); err != nil {
		log.Fatal(err)
	}
	var msg = make([]byte, 512)
	var n int
	if n, err = conn.Read(msg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received: %s.\n", msg[:n])
}

package main

import (
	"fmt"
	websocket "github.com/gorilla/websocket"
	"log"
	"net/url"
	"sync"
	"time"
)

func SendHeartBeat(ws *websocket.Conn) {
	interval := 25 * time.Second
	timer := time.NewTimer(25 * time.Second)
	defer timer.Stop()
	i := 0
	for {
		select {
		case <-timer.C:
			if i < 3 {
				timer.Reset(interval)
				i++
			}
			err := ws.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func main() {
	u := url.URL{
		Scheme: "ws",
		Host:   "localhost:9002",
		Path:   "/echo",
	}
	var conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	//这个设置为何没有起效？抓包看有收到Pong，但是没有执行到Handler内部 TODO 有空过下源码定位一下
	conn.SetPongHandler(func(appData string) error {
		fmt.Printf("Get Pong response: %s.\n", appData)
		err := conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		if err != nil {
			log.Fatal(err)
			return err
		}
		return nil
	})

	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()
		SendHeartBeat(conn)
	}()

	// send message
	sErr := conn.WriteMessage(websocket.TextMessage, []byte("{\"userId\":10001,\"code\":1001,\"data\":{}}"))
	if err != nil {
		log.Fatal(sErr)
	}
	// receive message
	_, message, rErr := conn.ReadMessage()
	if err != nil {
		log.Fatal(rErr)
	}
	fmt.Printf("Received: %s.\n", message)

	wg.Wait()
}

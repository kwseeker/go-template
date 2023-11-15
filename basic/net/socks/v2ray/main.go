package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

// 演示 v2ray 基于 socks Inbound 与 freedom Outbound 的代理流程
func main() {
	//1 配置加载
	configFile := "basic/net/socks/v2ray/config.json"
	config, err := LoadConfig(configFile)
	if err != nil {
		log.Panic("load config failed: ", err)
	}

	//2 服务器组件初始化
	server, err := NewServer(config)
	if err != nil {
		log.Panic("new server failed: ", err)
	}

	//3 服务器启动（即启动内部的 Inbound Server、Outbound Client）
	err = server.Start()
	defer server.Close()

	runtime.GC()
	//4 阻塞等待终止信号
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	<-osSignals
}

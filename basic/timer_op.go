package main

import (
	"log"
	"time"
)

// sleep.go runtimeTimer 和 time.go中的 timer 结构完全一样，可以互相转换
// 大概瞅了几秒源码，timer实现看上去类似Netty的时间轮
func main() {
	aTimer := time.NewTimer(3 * time.Second)
	defer aTimer.Stop()
	bTimer := time.NewTimer(10 * time.Second)
	defer bTimer.Stop()

	for {
		log.Print("loop ...")
		select {
		case <-aTimer.C:
			log.Print("aTimer timeout")
			aTimer.Reset(3 * time.Second)
		case <-bTimer.C:
			log.Print("bTimer timeout, exit")
			return
		}
	}
}

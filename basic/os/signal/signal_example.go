package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// 在终端中执行
// gops 查看进程ID，如果没装 gops: go install github.com/google/gops@latest
// kill -2 <pid>
func main() {
	done := make(chan bool, 1)

	s1 := make(chan os.Signal, 1)
	signal.Notify(s1, syscall.SIGINT, syscall.SIGTERM) //2 15

	go func() {
		<-s1
		fmt.Println("The program is going to exit ...")
		done <- true
	}()

	s2 := make(chan os.Signal, 1)
	signal.Notify(s2, syscall.SIGWINCH)

	go func() {
		for {
			<-s2
			fmt.Println("The terminal has been resized.")
		}
	}()

	<-done
}

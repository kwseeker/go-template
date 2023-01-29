package main

import (
	"fmt"
	"math/rand"
)

func Consumer(ch <-chan int, result chan<- int) {
	sum := 0
	for i := 0; i < 5; i++ {
		sum += <-ch
	}

	result <- sum
}

func Producer(ch chan<- int) {
	var num int
	for i := 0; i < 5; i++ {
		rand.Seed(20)
		num = rand.Intn(100)
		ch <- num
	}
}

func main() {
	ch := make(chan int)
	result := make(chan int)
	go Producer(ch)
	go Consumer(ch, result)

	fmt.Printf("result: %d\n", <-result)
}

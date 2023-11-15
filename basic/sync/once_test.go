package sync

import (
	"fmt"
	"sync"
	"testing"
)

// 多个协程中也只会执行一次
func TestOnce(t *testing.T) {
	//确保一个函数只执行一次
	var once sync.Once
	onceBody := func() {
		t.Logf("Only once")
	}
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			once.Do(onceBody)
			done <- true
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
	// Output:
	// Only once
}

// 单个协程中也只会执行一次
func TestOnce2(t *testing.T) {
	var count int
	increment := func() {
		count++
	}
	decrement := func() {
		count--
	}
	var once sync.Once
	once.Do(increment)
	once.Do(decrement)

	fmt.Printf("Count: %d\n", count)
}

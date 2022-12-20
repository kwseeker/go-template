package channel

import (
	"sync"
	"testing"
	"time"
)

func TestUnbuffered(t *testing.T) {
	//make chan 时不指定size就是创建不带缓冲的channel
	ch := make(chan string)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		ch <- `foo`            //写会阻塞，直到有读
		println(`write done!`) //注意比较这一句日志和无缓冲channel的输出位置
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Second)
		println(`Message: ` + <-ch)
	}()

	wg.Wait()
}

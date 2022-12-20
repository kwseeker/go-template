package channel

import (
	"sync"
	"testing"
	"time"
)

func TestBuffered(t *testing.T) {
	//make chan 时指定size就是创建带缓冲的channel, 这里1指可以缓冲写入1个字符串
	ch := make(chan string, 1)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		ch <- `foo`                 //写不会阻塞，除非缓冲被写满
		println(`write done!`)      //注意比较这一句日志和无缓冲channel的输出位置
		ch <- `bar`                 //缓冲已经写满，会阻塞
		println(`write more done!`) //注意比较这一句日志和无缓冲channel的输出位置
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Second)
		println(`Message: ` + <-ch)
		println(`Message: ` + <-ch)
	}()

	wg.Wait()
}

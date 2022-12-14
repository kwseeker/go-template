package sync

import (
	"runtime"
	"sync"
	"testing"
)

func HammerMutex(m *sync.Mutex, loops int, cdone chan bool, sum *int32) {
	for i := 0; i < loops; i++ {
		if i%3 == 0 {
			if m.TryLock() { //尝试获取锁，可能会失败，失败后不会执行+1，比如sum:9859, 说明有141次TryLock()失败
				*sum++
				m.Unlock()
			}
			continue
		}
		m.Lock() //阻塞等待获取锁
		*sum++
		m.Unlock()
	}
	cdone <- true
}

/*
go test -v -mutexprofile=mutex.prof -blockprofile=block.prof
go tool pprof mutex.prof
go tool pprof block.prof
*/
func TestMutex(t *testing.T) {
	//开启互斥锁分析，报告互斥锁的竞争情况,rate是采样率, 函数返回原rate
	if n := runtime.SetMutexProfileFraction(1); n != 0 {
		t.Logf("got mutexrate %d expected 0", n)
	}
	//最后关闭互斥锁分析
	defer runtime.SetMutexProfileFraction(0)
	runtime.SetBlockProfileRate(1)
	defer runtime.SetBlockProfileRate(0)

	m := new(sync.Mutex)
	sum := int32(0)

	m.Lock()
	t.Logf("在锁同步块内部，TryLock()会失败")
	if m.TryLock() {
		t.Fatalf("TryLock succeeded with mutex locked")
	}
	m.Unlock()
	t.Logf("已释放锁，TryLock()会成功")
	if !m.TryLock() {
		t.Fatalf("TryLock failed with mutex unlocked")
	}
	m.Unlock()

	//开启10个协程，每个协程不断轮询尝试获取锁并解锁（i是3的倍数 就 TryLock() Unlock(), 否则Lock() Unlock()）
	c := make(chan bool)
	rNum := 100
	for i := 0; i < rNum; i++ {
		go HammerMutex(m, 1000, c, &sum)
	}
	//等待所有协程执行完毕
	for i := 0; i < rNum; i++ {
		<-c
	}

	t.Logf("main routine done!, sum %d", sum)
	//_ = http.ListenAndServe(":6061", nil)
}

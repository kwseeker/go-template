package bytespool

import "sync"

const (
	numPools  = 4
	sizeMulti = 4
)

var (
	pool     [numPools]sync.Pool
	poolSize [numPools]int32
)

func createAllocFunc(size int32) func() interface{} {
	return func() interface{} {
		return make([]byte, size)
	}
}

// 创建了4个存储 []byte 的 sync.Pool，存储的 []byte size 分别为 2048 2048×4 2048×16 2048×64
func init() {
	size := int32(2048)
	for i := 0; i < numPools; i++ {
		pool[i] = sync.Pool{
			New: createAllocFunc(size),
		}
		poolSize[i] = size
		size *= sizeMulti
	}
}

// GetPool 返回生成 bytes 数组大小不小于 size 的对象池
func GetPool(size int32) *sync.Pool {
	for idx, ps := range poolSize {
		if size <= ps {
			return &pool[idx]
		}
	}
	return nil
}

// Alloc returns a byte slice with at least the given size. Minimum size of returned slice is 2048.
// 对象池有可用对象就返回，没有就 make() 创建，对象用完要放入对象池
func Alloc(size int32) []byte {
	pool := GetPool(size)
	if pool != nil {
		return pool.Get().([]byte)
	}
	return make([]byte, size)
}

// Free puts a byte slice into the internal pool.
// 将对象放入对象池
func Free(b []byte) {
	size := int32(cap(b))
	b = b[0:cap(b)]
	for i := numPools - 1; i >= 0; i-- {
		if size >= poolSize[i] {
			pool[i].Put(b)
			return
		}
	}
}

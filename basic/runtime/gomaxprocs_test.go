package runtime

import (
	"log"
	"runtime"
	"testing"
)

func TestGOMAXPROCS(t *testing.T) {
	n := runtime.GOMAXPROCS(4)
	defer runtime.GOMAXPROCS(n)

	log.Println("n=", n) //8, 本地机器是8核心
}

package gorm

import (
	"log"
	"testing"
)

func TestBatchQuery(t *testing.T) {
	BatchQuery()
}

type PairCounter struct {
	SelfCounter  int32
	LoverCounter int32
}

func TestBasic(t *testing.T) {
	var loverId int64
	log.Println(loverId == 0)

	//counterMap := make(map[int32]PairCounter)
	counterMap := make(map[int32]*PairCounter) //golang map元素无法取指，只能传指针才能修改值内部的字段
	if pairCounter, ok := counterMap[0]; ok {
		pairCounter.SelfCounter += 1
	} else {
		counterMap[0] = &PairCounter{0, 0}
	}

	counterMap[0].SelfCounter += 1
	log.Println(counterMap[0].SelfCounter, counterMap[0].LoverCounter)
}

package strconv

import (
	"strconv"
	"testing"
)

func TestStrconv(t *testing.T) {
	tests := []struct {
		s             string
		base          int
		targetBitSize int
		expectResult  uint64
	}{
		{"17", 8, 8, 15},
		{"66535", 10, 64, 66535},
	}

	for _, test := range tests {
		actual, err := strconv.ParseUint(test.s, test.base, test.targetBitSize)
		if err != nil {
			t.Fatal("parse failed")
		}
		t.Log("actual=", actual)
		if actual != test.expectResult {
			t.Error("expect:", test.expectResult, ", actual:", actual)
		}
	}
}

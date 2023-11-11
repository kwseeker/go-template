package io

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"
)

// 读取 json 文件并反序列化到 map[string]interface{}
func TestIo(t *testing.T) {
	content, _ := os.ReadFile("./test.json")
	reader := bytes.NewReader(content)
	n := make(map[string]interface{})

	buffer := bytes.NewBuffer(make([]byte, 0, 10240))
	teeReader := io.TeeReader(reader, buffer)
	_ = json.NewDecoder(teeReader).Decode(&n)

	if len(n) <= 0 {
		t.Error("Test Failed")
	}
}

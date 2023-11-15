package net

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// 端口以及端口范围解析
// 端口范围定义json格式支持 整数、整数的字符串、整数使用“-”连接的字符串

type Port uint16

func PortFromInt(val uint32) (Port, error) {
	if val > 65535 {
		return Port(0), errors.New(fmt.Sprint("invalid port range: ", val))
	}
	return Port(val), nil
}

func parseIntPort(data []byte) (Port, error) {
	var intPort uint32
	err := json.Unmarshal(data, &intPort)
	if err != nil {
		return Port(0), err
	}
	return PortFromInt(intPort)
}

func PortFromString(s string) (Port, error) {
	//将字符串按10进制，转成32位长的uint32类型
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return Port(0), errors.New(fmt.Sprint("invalid port range: ", s))
	}
	return PortFromInt(uint32(val))
}

func parseStringPort(s string) (Port, Port, error) {
	pair := strings.SplitN(s, "-", 2)
	if len(pair) == 0 {
		return Port(0), Port(0), errors.New(fmt.Sprint("invalid port range: ", s))
	}
	if len(pair) == 1 { //如"1080"
		port, err := PortFromString(pair[0])
		return port, port, err
	}
	fromPort, err := PortFromString(pair[0])
	if err != nil {
		return Port(0), Port(0), err
	}
	toPort, err := PortFromString(pair[1])
	if err != nil {
		return Port(0), Port(0), err
	}
	return fromPort, toPort, nil
}

// 解析类似 "1080" "1080-1090" 格式端口定义
func parseJSONStringPort(data []byte) (Port, Port, error) {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return Port(0), Port(0), err
	}
	return parseStringPort(s)
}

type PortRange struct {
	From uint32
	To   uint32
}

// UnmarshalJSON implements encoding/json.Unmarshaler.UnmarshalJSON
func (v *PortRange) UnmarshalJSON(data []byte) error {
	port, err := parseIntPort(data)
	if err == nil {
		v.From = uint32(port)
		v.To = uint32(port)
		return nil
	}

	from, to, err := parseJSONStringPort(data)
	if err == nil {
		v.From = uint32(from)
		v.To = uint32(to)
		if v.From > v.To {
			return errors.New(fmt.Sprint("invalid port range ", v.From, " -> ", v.To))
		}
		return nil
	}

	return errors.New(fmt.Sprint("invalid port range: ", string(data)))
}

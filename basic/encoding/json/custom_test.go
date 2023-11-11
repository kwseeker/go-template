// 定制序列化和反序列化规则
package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
)

type PortRange struct {
	From uint32
	To   uint32
}

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
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return Port(0), errors.New(fmt.Sprint("invalid port range: ", s))
	}
	return PortFromInt(uint32(val))
}

func parseStringPort(s string) (Port, Port, error) {
	if strings.HasPrefix(s, "env:") {
		s = s[4:]
		s = os.Getenv(s)
	}

	pair := strings.SplitN(s, "-", 2)
	if len(pair) == 0 {
		return Port(0), Port(0), errors.New(fmt.Sprint("invalid port range: ", s))
	}
	if len(pair) == 1 {
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

func parseJSONStringPort(data []byte) (Port, Port, error) {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return Port(0), Port(0), err
	}
	return parseStringPort(s)
}

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

func TestUnmarshalPortRange(t *testing.T) {
	var tests = []struct {
		jsonStr       string
		wantPortRange []uint32
	}{
		{"1080", []uint32{1080, 1080}},
		//注意这里必须把双引号用反斜杠带进去，否则实际就是这种格式 {"port": 1080-1090}，是有语法错误的，既不是数字也不是字符串
		//会报 “invalid character '-' after top-level value” 错误
		{"\"1080-1090\"", []uint32{1080, 1090}},
	}

	for _, test := range tests {
		pr := &PortRange{}
		err := json.Unmarshal([]byte(test.jsonStr), pr)
		if err != nil {
			t.Fatal("unmarshal failed: ", err)
		}
		if pr.From != test.wantPortRange[0] || pr.To != test.wantPortRange[1] {
			t.Errorf("jsonStr=%q, expect={From:%v,To:%v}, actual={From:%v,To:%v}",
				test.jsonStr, test.wantPortRange[0], test.wantPortRange[1], pr.From, pr.To)
		}
	}
}

func TestPort(t *testing.T) {
	port := Port(10)
	port2 := uint16(20)
	log.Println(port, ",", port2)
}

type Animal int

const (
	Unknown Animal = iota
	Gopher
	Zebra
)

// UnmarshalJSON 反序列化，这里动物只认 Gopher Zebra, 其他归为 Unknown
func (a *Animal) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	default:
		*a = Unknown
	case "gopher":
		*a = Gopher
	case "zebra":
		*a = Zebra
	}

	return nil
}

func (a Animal) MarshalJSON() ([]byte, error) {
	var s string
	switch a {
	default:
		s = "unknown"
	case Gopher:
		s = "gopher"
	case Zebra:
		s = "zebra"
	}

	return json.Marshal(s)
}

func TestCustomMarshalAndUnmarshal(t *testing.T) {
	//按定制规则进行序列化
	bs, err := json.Marshal(Gopher)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("gopher marshal: ", string(bs))

	//按定制规则进行反序列化，并计数
	blob := `["gopher","armadillo","zebra","unknown","gopher","bee","gopher","zebra"]`
	var zoo []Animal
	if err := json.Unmarshal([]byte(blob), &zoo); err != nil {
		log.Fatal(err)
	}

	census := make(map[Animal]int)
	for _, animal := range zoo {
		census[animal] += 1
	}

	log.Printf("Zoo Census:\n* Gophers: %d\n* Zebras:  %d\n* Unknown: %d\n",
		census[Gopher], census[Zebra], census[Unknown])
}

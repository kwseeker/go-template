package conf

import (
	"encoding/json"
	"errors"
	"fmt"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/common/net"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/common/protocol"
	"strings"
)

type StringList []string

func NewStringList(raw []string) *StringList {
	list := StringList(raw)
	return &list
}

func (v *StringList) Len() int {
	return len(*v)
}

func (v *StringList) UnmarshalJSON(data []byte) error {
	var strArr []string
	if err := json.Unmarshal(data, &strArr); err == nil {
		*v = *NewStringList(strArr)
		return nil
	}

	var rawStr string
	if err := json.Unmarshal(data, &rawStr); err == nil {
		strList := strings.Split(rawStr, ",")
		*v = *NewStringList(strList)
		return nil
	}
	return errors.New("unknown format of a string list: " + string(data))
}

type Address struct {
	net.Address
}

// UnmarshalJSON implements encoding/json.Unmarshaler.UnmarshalJSON
func (v *Address) UnmarshalJSON(data []byte) error {
	var rawStr string
	if err := json.Unmarshal(data, &rawStr); err != nil {
		return errors.New(fmt.Sprint("invalid address: ", string(data)))
	}
	v.Address = net.ParseAddress(rawStr)

	return nil
}

func (v *Address) Build() *protocol.IPOrDomain {
	return net.NewIPOrDomain(v.Address)
}

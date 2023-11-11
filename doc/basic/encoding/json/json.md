# Json

## 序列化与反序列化

`encoding/json` 包中实现了 对象 与 json 之间的序列化和反序列化。

这里主要研究**定制序列化和反序列化规则**。

常规情况，序列化和反序列化只需要将对象字段名和值和 json 的KV 一一对应赋值即可；

但是也存在一些非常规情况，需要自定义序列化和反序列化规则，比如下面的例子。

定制规则需要实现下面两个接口：

```go
type Marshaler interface {
    MarshalJSON() ([]byte, error)
}
type Unmarshaler interface {
    UnmarshalJSON([]byte) error
}
```

场景举例：比如看到网络代理开源项目 v2ray-core 的源码中，将端口 port 的 json 字符串反序列化成 PortRange 对象：

```json
{
	//端口配置支持三种格式，要求都可以反序列化成 PortRange 对象
    "port": "env:V2RAY_PORT_RANGE"
    //反序列化后 PortRange From=1080, To=1080
	"port": 1080
    //反序列化后 PortRange From=1080, To=1090
	"port": "1080-1090"
}
```

```go
type PortRange struct {
	From uint32
	To   uint32
}

// UnmarshalJSON 隐式实现了 encoding/json.Unmarshaler 接口
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
			return newError("invalid port range ", v.From, " -> ", v.To)
		}
		return nil
	}

	return newError("invalid port range: ", string(data))
}
```


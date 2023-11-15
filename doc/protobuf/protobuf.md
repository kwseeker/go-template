# Protobuf

[Protocol Buffer Basics: Go](https://protobuf.dev/getting-started/gotutorial/)

[proto3](https://protobuf.dev/programming-guides/proto3/)

[proto2](https://protobuf.dev/programming-guides/proto2/)

## 编写 .proto

### [定制 Options](https://protobuf.dev/programming-guides/proto2/#customoptions)

这部分主要在 proto2 中讲解。

其实是拓展官方的消息，因为选项是由Google/Protobuf/Deciptor.proto（例如FileOptions或FieldOptions）中定义的消息定义的。

Protocol Buffers 中，选项（Options）是键值对（Key-Value）的形式存储的。

当使用拓展选项时，选项名需要用圆括号括起来，以表明它是一个扩展，然后可以通过 Descriptor 读取选项的值。

自定义选项可以被定义为proto中的任意类型，如string, int32, enum, 甚至message。

案例：v2ray-core 中自定义 的选项：

```protobuf
extend google.protobuf.FieldOptions {
  FieldOpt field_opt = 50000;
}

message FieldOpt{
  repeated string any_wants = 1;
  repeated string allowed_values = 2;
  repeated string allowed_value_types = 3;

  // convert_time_read_file_into read a file into another field, and clear this field during input parsing
  string convert_time_read_file_into = 4;
  // forbidden marks a boolean to be inaccessible to user
  bool forbidden = 5;
  // convert_time_resource_loading read a file, and place its resource hash into another field
  string convert_time_resource_loading = 6;
  // convert_time_parse_ip parse a string ip address, and put its binary representation into another field
  string convert_time_parse_ip = 7;
}

message GeoIP {
  string country_code = 1;
  repeated CIDR cidr = 2;
  bool inverse_match = 3;

  // resource_hash instruct simplified config converter to load domain from geo file.
  bytes resource_hash = 4;
  string code = 5;

  string file_path = 68000 [(v2ray.core.common.protoext.field_opt).convert_time_resource_loading = "resource_hash"];
}

message GeoIPList {
  repeated GeoIP entry = 1;
}

message GeoSite {
  string country_code = 1;
  repeated Domain domain = 2;

  // resource_hash instruct simplified config converter to load domain from geo file.
  bytes resource_hash = 3;
  string code = 4;

  string file_path = 68000 [(v2ray.core.common.protoext.field_opt).convert_time_resource_loading = "resource_hash"];
}

message GeoSiteList {
  repeated GeoSite entry = 1;
}
```



## 编译 .proto

```shell
# protoc -I=$SRC_DIR --go_out=$DST_DIR $SRC_DIR/addressbook.proto
# -I 表示从哪个路径下查找依赖的 proto 文件，$SRC_DIR 搭配 import 语句组成最终的完整查找路径
# --go_out 表示生成文件路径，$DST_DIR 搭配 go_package 值组成最终生成文件的完整路径

# 不规范的操作：项目根目录执行（自动生成的go文件import会有问题）
protoc -I=./protobuf/proto --go_out=. protobuf/proto/addressbook.proto
protoc -I=./protobuf/proto --go_out=. protobuf/proto/*.proto
protoc -I=./protobuf/proto --go_out=. ./**/*.proto
protoc -I=. --go_out=. basic/net/socks/v2ray/common/protocol/*.proto
protoc -I=. --go_out=. basic/net/socks/**/*.proto

# 规范的操作：在$GOROOT/src下执行（这样自动生成的go文件没有任何问题可以直接运行）
# go.mod 模块命名： module kwseeker.top/kwseeker/go-template
# go_package 和 import 都使用标准的模块名加文件路径：
# 	option go_package = "kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/common/protocol";
# 	import "kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/common/protocol/address.proto";
protoc -I=. --go_out=. kwseeker.top/kwseeker/go-template/basic/net/socks/**/*.proto
```

关于 goland 中protobuf import 不识别问题：由于个人的项目都在`$GOROOT/src` 下 `kwseeker.top/kwseeker/project-name`，由 proto 文件查找规则，在 `Settings/Language&Frameworks/Protocol Buffers/Import Paths` 中添加 `$GOROOT/src` 即可。

## 写、读 Message

```go
//序列化
bytes, err := proto.Marshal(got)
cloned := &Person{}
//反序列化
err = proto.Unmarshal(bytes, cloned)
```


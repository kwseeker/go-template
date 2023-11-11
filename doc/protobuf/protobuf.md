# Protobuf

[Protocol Buffer Basics: Go](https://protobuf.dev/getting-started/gotutorial/)

## 编写 .proto

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


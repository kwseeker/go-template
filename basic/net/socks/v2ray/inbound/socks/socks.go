package socks

import (
	"context"
	"errors"
	"io"
	buf2 "kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/common/buf"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/infra/conf"
	"log"
	"net"
	"strings"
	"syscall"
	"time"
)

const (
	socks4Version = 0x04
	socks5Version = 0x05

	socks5AuthNotRequired      = 0x00
	socks5AuthGssAPI           = 0x01
	socks5AuthPassword         = 0x02
	socks5AuthNoMatchingMethod = 0xFF
)

type Server struct {
	config *ServerConfig
	//协议名
	Network string
	//IP 端口
	Address string
}

func NewServer(config *conf.Config) (*Server, error) {
	//TODO
	server := &Server{}
	return server, nil
}

type ServerSession struct {
	config *ServerConfig
	//address       net.Address
	//port          net.Port
	//clientAddress net.Address
}

// socks5 从第3个字节开始是 Methods 部分（认证方式列表），每个 Method 占一个字节，表示一种认证方式
// 如：	0x00 不需要认证
//		0x01 GSSAPI
//		0x02 用户名、密码认证
// 		......
// reader writer 即网络连接的读写接口
func (s *ServerSession) auth5(nMethod byte, reader io.Reader, writer io.Writer) (username string, err error) {
	buffer := buf2.StackNew()
	defer buffer.Release()

	//读取所有认证方法到 Buffer
	if _, err = buffer.ReadFullFrom(reader, int32(nMethod)); err != nil {
		return "", errors.New("failed to read auth methods")
	}

	//var expectedAuth byte = socks5AuthNotRequired
	//if s.config.AuthType == AuthType_PASSWORD {
	//	expectedAuth = socks5AuthPassword
	//}

	//if !hasAuthMethod(expectedAuth, buffer.BytesRange(0, int32(nMethod))) {
	//	writeSocks5AuthenticationResponse(writer, socks5Version, authNoMatchingMethod)
	//	return "", errors.New("no matching auth method")
	//}
	//
	//if err := writeSocks5AuthenticationResponse(writer, socks5Version, expectedAuth); err != nil {
	//	return "", errors.New("failed to write auth response")
	//}
	//
	//// 0x02 用户名、密码认证
	//if expectedAuth == socks5AuthPassword {
	//	username, password, err := ReadUsernamePassword(reader)
	//	if err != nil {
	//		return "", errors.New("failed to read username and password for authentication")
	//	}
	//
	//	if !s.config.HasAccount(username, password) {
	//		writeSocks5AuthenticationResponse(writer, 0x01, 0xFF)
	//		return "", errors.New("invalid username or password")
	//	}
	//
	//	if err := writeSocks5AuthenticationResponse(writer, 0x01, 0x00); err != nil {
	//		return "", errors.New("failed to write auth response")
	//	}
	//	return username, nil
	//}

	return "", nil
}

//func (s *ServerSession) handshake5(nMethods byte, reader io.Reader, writer io.Writer) (*protocol.RequestHeader, error) {
//	var username string
//	var err error
//	if username, err = s.auth5(nMethods, reader, writer); err != nil {
//		return nil, err
//	}
//
//}
//
//// handshake socks 握手
//// 包括认证
//// 用 Java 的习惯表示就是 reader -> Request writer -> Response
//func (s *ServerSession) handshake(reader io.Reader, writer io.Writer) (*protocol.RequestHeader, error) {
//	//读取连接请求内容，判断socks版本等信息
//	buffer := buf2.StackNew() //从对象池获取bytes,封成Buffer
//	defer buffer.Release()
//	if _, err := buffer.ReadFullFrom(reader, 2); err != nil {
//		buffer.Release()
//		return nil, errors.New("insufficient header")
//	}
//
//	// https://zh.wikipedia.org/zh-hans/SOCKS
//	socksVersion := buffer.Byte(0) //第一个字节是 socks 版本
//
//	switch socksVersion {
//	case socks4Version:
//		return nil, errors.New(fmt.Sprint("not support Socks version: ", socksVersion))
//	case socks5Version:
//		socksNMethods := buffer.Byte(1) //第二个字节是 METHODS 部分的长度
//		return s.handshake5(socksNMethods, reader, writer)
//	default:
//		return nil, errors.New(fmt.Sprint("unknown Socks version: ", socksVersion))
//	}
//}

func getControlFunc(ctx context.Context) func(network, address string, c syscall.RawConn) error {
	return func(network, address string, c syscall.RawConn) error {
		log.Println("do nothing now!")
		return nil
	}
}

// socks5 服务端需要在建立与客户端的TCP连接之后，解析请求内容确认socks协议版本及认证方式
// 然后根据 socks 版本进行握手
func handleProxy(conn net.Conn) {
	//session := &ServerSession{
	//	config: s.config,
	//}
	//
	////读取请求内容并解析
	//reader := &buf2.BufferedReader{
	//	Reader: buf2.NewReader(conn),
	//}
	//
	//session.handshake(reader, conn)

	_ = conn.Close()
}

// 监听接入请求
func (s *Server) keepAccepting(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			errStr := err.Error()
			log.Println("err = ", errStr)
			if strings.Contains(errStr, "closed") {
				break
			}
			continue
		}

		go handleProxy(conn)
	}
}

func (s *Server) Start(ctx context.Context) error {
	lc := net.ListenConfig{
		//创建网络连接后绑定之前调用此函数
		Control:   getControlFunc(ctx),
		KeepAlive: time.Duration(-1),
	}

	l, err := lc.Listen(ctx, s.Network, s.Address)
	if err != nil {
		return err
	}

	go s.keepAccepting(l)

	return nil
}

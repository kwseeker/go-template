package conf

import "encoding/json"

type TCPConfig struct {
	HeaderConfig        json.RawMessage `json:"header"`
	AcceptProxyProtocol bool            `json:"acceptProxyProtocol"`
}

type HTTPConfig struct {
	Host    *[]string            `json:"host"`
	Path    string               `json:"path"`
	Method  string               `json:"method"`
	Headers map[string]*[]string `json:"headers"`
}

type TransportConfig struct {
	TCPConfig  *TCPConfig  `json:"tcpSettings"`
	HTTPConfig *HTTPConfig `json:"httpSettings"`
	//KCPConfig  *KCPConfig          `json:"kcpSettings"`
	//WSConfig   *WebSocketConfig    `json:"wsSettings"`
	//DSConfig   *DomainSocketConfig `json:"dsSettings"`
	//QUICConfig *QUICConfig         `json:"quicSettings"`
	//GunConfig  *GunConfig          `json:"gunSettings"`
	//GRPCConfig *GunConfig          `json:"grpcSettings"`
}

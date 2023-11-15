package conf

import (
	"encoding/json"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/common/net"
)

type InboundDetourAllocationConfig struct {
	Strategy    string  `json:"strategy"`
	Concurrency *uint32 `json:"concurrency"`
	RefreshMin  *uint32 `json:"refresh"`
}

type InboundDetourConfig struct {
	Protocol   string                         `json:"protocol"`
	PortRange  *net.PortRange                 `json:"port"`
	ListenOn   *Address                       `json:"listen"`
	Settings   *json.RawMessage               `json:"settings"`
	Tag        string                         `json:"tag"`
	Allocation *InboundDetourAllocationConfig `json:"allocate"`
	//StreamSetting  *StreamConfig                  `json:"streamSettings"`
	DomainOverride *[]string `json:"domainOverride"`
	//SniffingConfig *sniffer.SniffingConfig `json:"sniffing"`
}

type OutboundDetourConfig struct {
	Protocol       string           `json:"protocol"`
	SendThrough    Address          `json:"sendThrough"`
	Tag            string           `json:"tag"`
	Settings       *json.RawMessage `json:"settings"`
	DomainStrategy string           `json:"domainStrategy"`
}

// Config v4 https://www.v2fly.org/config/overview.html
// 现在是配置文件先转成v4的配置类，再转成v5的配置类
type Config struct {
	//DNSConfig       *dns.DNSConfig         `json:"dns"`       //DNS服务, 对目标地址（域名）进行 DNS 解析，同时为 IP 路由规则匹配提供判断依据
	RouterConfig    *RouterConfig          `json:"routing"`   //路由配置，当有多个Outbound用于选择从哪个Outbound出站
	Policy          *PolicyConfig          `json:"policy"`    //本地策略，配置一些用户相关权限，可以配置多组，比如：超时设置，https://www.v2fly.org/config/policy.html#levelpolicyobject
	InboundConfigs  []InboundDetourConfig  `json:"inbounds"`  //入站规则
	OutboundConfigs []OutboundDetourConfig `json:"outbounds"` //出站规则
	Transport       *TransportConfig       `json:"transport"` //当前 V2Ray 节点和其它节点对接的方式，即 Outbound 向其他节点转发的协议配置
	//LogConfig        *log.LogConfig          `json:"log"`
	//API              *APIConfig              `json:"api"`     //API接口，v2ray有好几个客户端，都需要通过这些API修改 v2ray-core 配置
	//Stats            *StatsConfig            `json:"stats"`	//运行状态统计
	//Reverse          *ReverseConfig          `json:"reverse"`	//反向代理
	//FakeDNS          *dns.FakeDNSConfig      `json:"fakeDns"` //本质是记忆“IP-域名”映射，像是DNS的缓存服务
	//BrowserForwarder *BrowserForwarderConfig `json:"browserForwarder"` 	//浏览器转发
	//Observatory      *ObservatoryConfig      `json:"observatory"`			//连接观测组件，可以将观测结果反馈给其他组件，比如自动优化负载均衡
	//BurstObservatory *BurstObservatoryConfig `json:"burstObservatory"`
	//MultiObservatory *MultiObservatoryConfig `json:"multiObservatory"`

	//Services map[string]*json.RawMessage `json:"services"`
}

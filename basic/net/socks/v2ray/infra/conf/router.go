package conf

import (
	"context"
	"encoding/json"
	"strings"
)

type DomainStrategy int32

const (
	DomainStrategyAsIs DomainStrategy = iota
	DomainStrategyUseIp
	DomainStrategyIpIfNonMatch
	DomainStrategyIpOnDemand
)

var (
	DomainStrategyMap = map[string]DomainStrategy{
		"AsIs":         DomainStrategyAsIs,
		"UseIP":        DomainStrategyUseIp,
		"IPIfNonMatch": DomainStrategyIpIfNonMatch,
		"IPOnDemand":   DomainStrategyIpOnDemand,
	}
)

//type StrategyConfig struct {
//	Type     string           `json:"type"`
//	Settings *json.RawMessage `json:"settings"`
//}

//type BalancingRule struct {
//	Tag         string         `json:"tag"`
//	Selectors   []string       `json:"selector"`
//	Strategy    StrategyConfig `json:"strategy"`
//	FallbackTag string         `json:"fallbackTag"`
//}

// RouterConfig
//{
//	"domainStrategy": "IPOnDemand",
//	"rules":[{
//		"type": "field",
//		"domains": ["geosite:cn"],
//		"ip": ["geoip:cn"],
//		"outboundTag": "direct"
//	},{
//		"type": "field",
//		"domains": [
//			"geosite:google",
//			"geosite:tld-!cn",
//			"geolocation-!cn"
//		],
//		"ip": ["geoip:!cn"],
//		"outboundTag": "vps-http"
//	}]
//},
type RouterConfig struct { // nolint: revive
	DomainStrategy *string           `json:"domainStrategy"`
	RuleList       []json.RawMessage `json:"rules"`
	//Balancers      []*BalancingRule  `json:"balancers"`
	//DomainMatcher  string            `json:"domainMatcher"`
	cfgctx context.Context
}

func (c *RouterConfig) GetDomainStrategy() DomainStrategy {
	ds := ""
	if c.DomainStrategy != nil {
		ds = *c.DomainStrategy
	}

	switch strings.ToLower(ds) {
	case "alwaysip", "always_ip", "always-ip":
		return DomainStrategyUseIp
	case "ipifnonmatch", "ip_if_non_match", "ip-if-non-match":
		return DomainStrategyIpIfNonMatch
	case "ipondemand", "ip_on_demand", "ip-on-demand":
		return DomainStrategyIpOnDemand
	default:
		return DomainStrategyAsIs
	}
}

type RouterRule struct {
	Type        string `json:"type"`
	OutboundTag string `json:"outboundTag"`
	//BalancerTag string `json:"balancerTag"`
	//DomainMatcher string `json:"domainMatcher"`
}

type RawFieldRule struct {
	RouterRule
	Domains *StringList `json:"domains"`
	IP      *StringList `json:"ip"`
}

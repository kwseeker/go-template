package router

import (
	"encoding/json"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/common/errors"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/infra/conf"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/infra/geodata"
	"strings"
)

// RoutingRule 对应配置中的 "routing.rule" https://www.v2fly.org/config/routing.html
// {
//   "domains": [
//	 	"geosite:cn"
//   ],
//   "ip": [
//	 	"geoip:cn"
//   ],
//   "outboundTag": "direct",
//   "type": "field"
// }
type RoutingRule struct {
	//outboundTag
	TargetTag string
	//domains
	Domain []*conf.Domain
	//ip
	GeoIP []*conf.GeoIP
}

func parseDomainRule(domain string) ([]*conf.Domain, error) {
	geoLoader := geodata.GetGeoLoader()

	if strings.HasPrefix(domain, "geosite:") {
		list := domain[8:]
		if len(list) == 0 {
			return nil, errors.NewError("empty list name in rule: ", domain)
		}
		domains, err := geoLoader.LoadGeoSite(list)
		if err != nil {
			return nil, errors.NewError("failed to load geosite: ", list)
		}

		return domains, nil
	}

	domainRule := new(conf.Domain)
	domainRule.Type = conf.Domain_Plain
	domainRule.Value = domain

	return []*conf.Domain{domainRule}, nil
}

func toCidrList(ips conf.StringList) ([]*conf.GeoIP, error) {
	geoLoader := geodata.GetGeoLoader()

	var geoipList []*conf.GeoIP
	for _, ip := range ips {
		if strings.HasPrefix(ip, "geoip:") {
			country := ip[6:]
			isReverseMatch := false
			if strings.HasPrefix(ip, "geoip:!") {
				country = ip[7:]
				isReverseMatch = true
			}
			if len(country) == 0 {
				return nil, errors.NewError("empty country name in rule")
			}
			geoip, err := geoLoader.LoadGeoIP(country)
			if err != nil {
				return nil, errors.NewError("failed to load geoip: ", country)
			}

			geoipList = append(geoipList, &conf.GeoIP{
				CountryCode:  strings.ToUpper(country),
				Cidr:         geoip,
				InverseMatch: isReverseMatch,
			})
		}
	}

	return geoipList, nil
}

func newRoutingRule(rawRule json.RawMessage) (*RoutingRule, error) {
	//和配置文件格式对应的路由规则
	rawFieldRule := new(conf.RawFieldRule)
	err := json.Unmarshal(rawRule, rawFieldRule)
	if err != nil {
		return nil, err
	}

	//conf.RawFieldRule 转成 RoutingRule
	//比如 geosite:cn geoip:cn 这种需要转成geo数据库中的值
	rule := new(RoutingRule)
	//RoutingRule.TargetTag
	switch {
	case len(rawFieldRule.OutboundTag) > 0:
		rule.TargetTag = rawFieldRule.OutboundTag
	default:
		return nil, errors.NewError("neither outboundTag nor balancerTag is specified in routing rule")
	}
	//RoutingRule.Domain RoutingRule.GeoIP
	//Domain GeoIP 规则解析， geosite.dat geoip.dat 保存的都是 protobuf 序列化后的数据， 原数据参考官方文档链接
	if rawFieldRule.Domains != nil {
		for _, domain := range *rawFieldRule.Domains {
			rules, err := parseDomainRule(domain)
			if err != nil {
				return nil, errors.NewError("failed to parse domain rule: ", domain)
			}
			rule.Domain = append(rule.Domain, rules...)
		}
	}
	if rawFieldRule.IP != nil {
		geoipList, err := toCidrList(*rawFieldRule.IP)
		if err != nil {
			return nil, err
		}
		rule.GeoIP = geoipList
	}

	return rule, nil
}

// Config v5
type Config struct {
	DomainStrategy conf.DomainStrategy
	Rule           []*RoutingRule
}

// NewConfig 从 conf.RouterConfig 创建 router.Config
func NewConfig(rawConfig *conf.RouterConfig) (*Config, error) {
	config := new(Config)
	config.DomainStrategy = rawConfig.GetDomainStrategy()
	//[]json.RawMessage
	for _, rawRule := range rawConfig.RuleList {
		rule, err := newRoutingRule(rawRule)
		if err != nil {
			return nil, err
		}
		//if rule.DomainMatcher == "" {
		//	rule.DomainMatcher = c.DomainMatcher
		//}
		config.Rule = append(config.Rule, rule)
	}

	return config, nil
}

//{
//    "domainStrategy": "IPOnDemand",
//    "rules": [
//      {
//        "domains": [
//          "geosite:cn"
//        ],
//        "ip": [
//          "geoip:cn"
//        ],
//        "outboundTag": "direct",
//        "type": "field"
//      },
//      {
//        "domains": [
//          "geosite:google",
//          "geosite:tld-!cn",
//          "geolocation-!cn"
//        ],
//        "ip": [
//          "geoip:!cn"
//        ],
//        "outboundTag": "vps-http",
//        "type": "field"
//      }
//    ]
//  }

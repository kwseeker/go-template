package router

import (
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/features/outbound"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/features/routing"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/infra/conf"
	"sync"
)

//这里还是以v4版本的配置为主

type Balancer struct {
	selectors   []string
	strategy    BalancingStrategy
	ohm         outbound.Manager
	fallbackTag string

	override override
}

type Condition interface {
	Apply(ctx routing.Context) bool
}

// Rule 是RoutingRule经处理后真正参与逻辑的路由规则
type Rule struct {
	// outboundTag
	Tag string
	// v4负载均衡配置方式
	Balancer *Balancer
	// 匹配规则
	Condition Condition
}

type BalancingStrategy interface {
	PickOutbound([]string) string
}

type overrideSettings struct {
	target string
}

type override struct {
	access   sync.RWMutex
	settings overrideSettings
}

type Router struct {
	domainStrategy conf.DomainStrategy
	//路由规则，从上到下依次匹配
	rules []*Rule
	//这个是v5负载均衡配置的方式
	//balancers map[string]*Balancer
	//暂时不需要这个
	//dns            dns.Client
}

func New(config *conf.RouterConfig) (*Router, error) {
	routerConfig, err := NewConfig(config)
	if err != nil {
		return nil, err
	}

	r := new(Router)
	r.domainStrategy = r.GetDomainStrategy()
	//r.dns = d
	//r.balancers = make(map[string]*Balancer, len(config.BalancingRule))
	r.rules = make([]*Rule, 0, len(routerConfig.Rule))
	//TODO
	//for _, rule := range routerConfig.Rule {
	//	cond, err := rule.BuildCondition()
	//	if err != nil {
	//		return nil, err
	//	}
	//	rr := &Rule{
	//		Condition: cond,
	//		Tag:       rule.GetTag(),
	//	}
	//	r.rules = append(r.rules, rr)
	//}
	return r, nil
}

// Start implements common.Runnable.
func (r *Router) Start() error {
	return nil
}

// Close implements common.Closable.
func (r *Router) Close() error {
	return nil
}

// Type implements common.HasType.
func (*Router) Type() interface{} {
	return routing.RouterType()
}

func (r *Router) GetDomainStrategy() conf.DomainStrategy {
	return r.domainStrategy
}

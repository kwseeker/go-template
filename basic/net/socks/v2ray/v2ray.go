package main

import (
	"context"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/features"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/inbound"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/infra/conf"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/outbound"
	"kwseeker.top/kwseeker/go-template/basic/net/socks/v2ray/router"
	"log"
	"sync"
)

type Instance struct {
	ctx      context.Context
	access   sync.Mutex
	running  bool               //运行状态
	features []features.Feature //功能组件
	//featureResolutions []resolution
	//env                environment.RootEnvironment
}

func NewServer(config *conf.Config) (*Instance, error) {
	server := &Instance{ctx: context.Background()}

	done, err := initInstanceWithConfig(config, server)
	if done {
		return nil, err
	}

	return server, nil
}

// 组件实例化
//v2ray-core 源码在初始化阶段在 type.go 中为各个组件注册了一组构造器函数（内部也是&Xxx{}），在这个方法中调用，这里简略点直接创建
func initInstanceWithConfig(config *conf.Config, server *Instance) (bool, error) {
	//config.InboundConfigs -> inbound.Manager
	inboundManager, err := inbound.New(server.ctx)
	if err != nil {
		return false, err
	}
	if err := server.AddFeature(inboundManager); err != nil {
		return false, err
	}
	//config.OutboundConfigs -> outbound.Manager
	outboundManager, err := outbound.New()
	if err != nil {
		return false, err
	}
	if err := server.AddFeature(outboundManager); err != nil {
		return false, err
	}
	//config.RouterConfig -> router.Router
	r, err := router.New(config.RouterConfig)
	if err != nil {
		return false, err
	}
	if err := server.AddFeature(r); err != nil {
		return false, err
	}
	//config.Policy
	//config.Transport

	return true, nil
}

func (s *Instance) AddFeature(feature features.Feature) error {
	s.features = append(s.features, feature)
	if s.running {
		if err := feature.Start(); err != nil {
			log.Println("failed to start feature")
		}
		return nil
	}

	return nil
}

func (s *Instance) Start() error {

	return nil
}

func (s *Instance) Close() {

}

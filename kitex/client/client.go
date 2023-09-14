package client

import (
	"context"
	"net"
	"strings"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"

	"github.com/xh-polaris/gopkg/kitex/middleware"
	"github.com/xh-polaris/gopkg/util/log"
)

const (
	EnvHeader     = "X-Xh-Env"
	LaneHeader    = "X-Xh-Lane"
	magicEndpoint = "magic-host:magic-port"
)

func NewClient[C any](fromName, toName string, fn func(fromName string, opts ...client.Option) (C, error), directEndpoints ...string) C {
	cli, err := fn(
		fromName,
		client.WithHostPorts(func() []string {
			if len(directEndpoints) != 0 {
				return directEndpoints
			}
			return []string{magicEndpoint}
		}()...),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: fromName}),
		client.WithInstanceMW(middleware.LogMiddleware(toName)),
		client.WithLoadBalancer(&LoadBalancer{ServiceName: strings.ReplaceAll(toName, ".", "-")}),
	)
	if err != nil {
		log.Error("[NewClient], err=%v", err)
	}
	return cli
}

type LoadBalancer struct {
	ServiceName string
}

func (b *LoadBalancer) GetPicker(result discovery.Result) loadbalance.Picker {
	return &Picker{
		ServiceName: b.ServiceName,
		Instances:   result.Instances,
	}
}

func (b *LoadBalancer) Name() string {
	return "magic-name"
}

type Picker struct {
	ServiceName string
	Instances   []discovery.Instance
}

func (p *Picker) Next(ctx context.Context, _ interface{}) discovery.Instance {
	if len(p.Instances) != 0 && p.Instances[0].Address().String() != magicEndpoint {
		return p.Instances[0]
	}

	var host = p.ServiceName + ".xh-polaris"

	// 选择基准环境
	env, ok := metainfo.GetPersistentValue(ctx, EnvHeader)
	if ok && env == "test" {
		host += "-test"
	}

	// 检查泳道是否部署该服务
	lane, ok := metainfo.GetPersistentValue(ctx, LaneHeader)
	if ok && lane != "" {
		addr, err := net.ResolveTCPAddr("tcp", host+"-"+lane+".svc.cluster.local:8080")
		if err == nil {
			return &Instance{addr: addr}
		}
	}

	addr, err := net.ResolveTCPAddr("tcp", host+".svc.cluster.local:8080")
	if err == nil {
		return &Instance{addr: addr}
	}
	return nil
}

type Instance struct {
	addr net.Addr
}

func (i *Instance) Address() net.Addr {
	return i.addr
}

func (i *Instance) Weight() int {
	return 0
}

func (i *Instance) Tag(_ string) (value string, exist bool) {
	return "", false
}

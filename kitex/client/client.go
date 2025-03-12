package client

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/discovery"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/remote/trans/nphttp2/metadata"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/cloudwego/kitex/transport"
	prometheus "github.com/kitex-contrib/monitor-prometheus"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"k8s.io/utils/env"

	"github.com/xh-polaris/gopkg/consts"
	"github.com/xh-polaris/gopkg/kitex/middleware"
	"github.com/xh-polaris/gopkg/util/log"
)

var tracer = prometheus.NewClientTracer(":9091", "/client/metrics")

func NewClient[C any](fromName, toName string, fn func(string, ...client.Option) (C, error), opts ...client.Option) C {
	opts = append([]client.Option{
		WithHostPorts(fmt.Sprintf("%s:8080", strings.ReplaceAll(toName, ".", "-"))),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithTracer(tracer),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: fromName}),
		client.WithInstanceMW(middleware.LogMiddleware(toName)),
		client.WithInstanceMW(middleware.EnvironmentMiddleware),
		client.WithTransportProtocol(transport.GRPC),
		client.WithMetaHandler(transmeta.ClientHTTP2Handler)},
		opts...,
	)
	if val, err := env.GetBool("IGNORE_MESH", false); err != nil && val {
		opts = append(opts, client.WithLoadBalancer(&LoadBalancer{ServiceName: strings.ReplaceAll(toName, ".", "-")}))
	}
	cli, err := fn(
		fromName,
		opts...,
	)
	if err != nil {
		log.Error("[NewClient], toName=%s, err=%v", toName, err)
	}
	return cli
}

func WithHostPorts(hostports ...string) client.Option {
	return client.Option{F: func(o *client.Options, di *utils.Slice) {
		di.Push(fmt.Sprintf("WithHostPorts(%v)", hostports))

		var ins []discovery.Instance
		for _, hp := range hostports {
			if _, err := net.ResolveTCPAddr("tcp", hp); err == nil {
				ins = append(ins, discovery.NewInstance("tcp", hp, discovery.DefaultWeight, nil))
				continue
			}
			panic(fmt.Errorf("WithHostPorts: invalid '%s'", hp))
		}

		if len(ins) == 0 {
			panic("WithHostPorts() requires at least one argument")
		}

		o.Targets = strings.Join(hostports, ",")
		o.Resolver = &discovery.SynthesizedResolver{
			ResolveFunc: func(ctx context.Context, key string) (discovery.Result, error) {
				return discovery.Result{
					Cacheable: true,
					CacheKey:  "fixed",
					Instances: ins,
				}, nil
			},
			NameFunc: func() string { return o.Targets },
			TargetFunc: func(ctx context.Context, target rpcinfo.EndpointInfo) string {
				return o.Targets
			},
		}
	}}
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
	var host = p.ServiceName + ".xh-polaris"

	// 选择基准环境
	env, ok := metainfo.GetPersistentValue(ctx, consts.EnvHeader)
	if !ok {
		var md metadata.MD
		md, ok = metadata.FromIncomingContext(ctx)
		if ok && len(md[consts.EnvHeader]) > 0 {
			env = md[consts.EnvHeader][0]
		}
	}
	if ok && env == "test" {
		host += "-test"
	}

	// 检查泳道是否部署该服务
	lane, ok := metainfo.GetPersistentValue(ctx, consts.LaneHeader)
	if !ok {
		var md metadata.MD
		md, ok = metadata.FromIncomingContext(ctx)
		if ok && len(md[consts.LaneHeader]) > 0 {
			lane = md[consts.LaneHeader][0]
		}
	}
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

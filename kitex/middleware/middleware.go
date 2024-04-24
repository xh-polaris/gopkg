package middleware

import (
	"context"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/remote/trans/nphttp2/metadata"

	"github.com/xh-polaris/gopkg/consts"
	"github.com/xh-polaris/gopkg/util"
	"github.com/xh-polaris/gopkg/util/log"
)

var (
	LogMiddleware = func(name string) endpoint.Middleware {
		return func(next endpoint.Endpoint) endpoint.Endpoint {
			return func(ctx context.Context, req, resp interface{}) error {
				err := next(ctx, req, resp)
				log.CtxInfo(ctx, "[%s RPC Request] req=%s, resp=%s, err=%v", name, util.JSONF(req), util.JSONF(resp), err)
				return err
			}
		}
	}
	EnvironmentMiddleware = func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req, resp interface{}) error {
			md, _ := metadata.FromIncomingContext(ctx)
			if env := md.Get(consts.EnvHeader); len(env) != 0 && env[0] != "" {
				ctx = metadata.AppendToOutgoingContext(ctx, consts.EnvHeader, env[0])
			}
			return next(ctx, req, resp)
		}
	}
)

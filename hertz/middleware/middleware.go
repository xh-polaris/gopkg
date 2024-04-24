package middleware

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"google.golang.org/grpc/metadata"

	"github.com/xh-polaris/gopkg/consts"
)

var (
	EnvironmentMiddleware = func(ctx context.Context, c *app.RequestContext) {
		if env := c.Request.Header.Get(consts.EnvHeader); env != "" {
			if lane := c.Request.Header.Get(consts.LaneHeader); lane != "" {
				env = fmt.Sprintf("%s_%s", env, lane)
			}
			ctx = metadata.AppendToOutgoingContext(ctx, consts.EnvHeader, env)
		}
		c.Next(ctx)
	}
)

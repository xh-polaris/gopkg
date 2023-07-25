package middleware

import (
	"context"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/hertz/pkg/app"
)

var (
	EnvironmentMiddleware = func(ctx context.Context, c *app.RequestContext) {
		if env := c.Request.Header.Get("X-Xh-Env"); env != "" {
			ctx = metainfo.WithPersistentValue(ctx, "X-Xh-Env", env)
		}
		if lane := c.Request.Header.Get("X-Xh-Lane"); lane != "" {
			ctx = metainfo.WithPersistentValue(ctx, "X-Xh-Lane", lane)
		}
		c.Next(ctx)
	}
)

package middleware

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// AuthMiddleware enforces authentication for non-public endpoints.
func AuthMiddleware(
	publicPaths map[string]struct{},
	validator func(operation string, header transport.Header) error,
) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			transportInfo, ok := transport.FromServerContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			if _, allowed := publicPaths[transportInfo.Operation()]; allowed {
				return handler(ctx, req)
			}

			header := transportInfo.RequestHeader()
			if header.Get("Authorization") == "" && header.Get("X-User") == "" {
				if err := validator(transportInfo.Operation(), header); err != nil {
					return nil, errors.Unauthorized("UNAUTHORIZED", err.Error())
				}
			}

			return handler(ctx, req)
		}
	}
}

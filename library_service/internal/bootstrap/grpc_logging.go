package bootstrap

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func unaryLoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)

		st := status.Convert(err)
		code := st.Code()
		if err == nil {
			code = codes.OK
		}

		attrs := []any{
			"method", info.FullMethod,
			"code", code.String(),
			"duration_ms", time.Since(start).Milliseconds(),
		}
		if err != nil {
			slog.Warn("grpc request failed", attrs...)
		} else {
			slog.Info("grpc request", attrs...)
		}

		return resp, err
	}
}


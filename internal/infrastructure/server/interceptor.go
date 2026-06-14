package server

import (
	"context"
	"runtime/debug"

	"github.com/Junaidmdv/goalcircle-team_service/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RecoveryInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp any, err error) {

		defer func() {
			if r := recover(); r != nil {
				logger.Error(
					"panic recovered",
					"method", info.FullMethod,
					"panic", r,
					"stack", string(debug.Stack()),
				)
				err = status.Error(
					codes.Internal,
					"internal server error",
				)

			}
		}()

		return handler(ctx, req)
	}
}

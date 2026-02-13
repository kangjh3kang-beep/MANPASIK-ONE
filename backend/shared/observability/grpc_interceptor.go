package observability

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// UnaryServerInterceptor returns a gRPC unary interceptor that records metrics
func UnaryServerInterceptor(m *Metrics) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		statusCode := 200
		if err != nil {
			statusCode = 500
		}
		m.RecordRequest(info.FullMethod, duration, statusCode)

		return resp, err
	}
}

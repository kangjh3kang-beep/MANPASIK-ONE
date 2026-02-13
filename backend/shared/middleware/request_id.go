package middleware

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const RequestIDKey contextKey = "request_id"
const RequestIDHeader = "x-request-id"

// RequestIDInterceptor adds a unique request ID to each request
func RequestIDInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		requestID := ""

		// Try to get from incoming metadata
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if ids := md.Get(RequestIDHeader); len(ids) > 0 {
				requestID = ids[0]
			}
		}

		// Generate if not present
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Add to context
		ctx = context.WithValue(ctx, RequestIDKey, requestID)

		// Add to outgoing metadata for propagation
		ctx = metadata.AppendToOutgoingContext(ctx, RequestIDHeader, requestID)

		return handler(ctx, req)
	}
}

// RequestIDFromContext extracts request ID from context
func RequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

package middleware

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const tracerName = "manpasik.middleware.tracing"

// UnaryTracingInterceptor returns a gRPC unary server interceptor that
// creates a span for each RPC call, propagates trace context, and records
// status/error attributes.
func UnaryTracingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		tracer := otel.Tracer(tracerName)

		ctx = extractTraceContext(ctx)

		ctx, span := tracer.Start(ctx, info.FullMethod,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("rpc.system", "grpc"),
				attribute.String("rpc.method", info.FullMethod),
			),
		)
		defer span.End()

		// Propagate request ID if available
		if reqID := RequestIDFromContext(ctx); reqID != "" {
			span.SetAttributes(attribute.String("request.id", reqID))
		}

		// Propagate user ID if available
		if userID, ok := UserIDFromContext(ctx); ok {
			span.SetAttributes(attribute.String("user.id", userID))
		}

		resp, err := handler(ctx, req)
		if err != nil {
			st, _ := status.FromError(err)
			span.SetStatus(codes.Error, st.Message())
			span.SetAttributes(
				attribute.String("rpc.grpc.status_code", st.Code().String()),
			)
			span.RecordError(err)
		} else {
			span.SetStatus(codes.Ok, "")
			span.SetAttributes(
				attribute.String("rpc.grpc.status_code", grpcCodes.OK.String()),
			)
		}

		return resp, err
	}
}

// StreamTracingInterceptor returns a gRPC stream server interceptor that
// creates a span for the stream lifecycle.
func StreamTracingInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		tracer := otel.Tracer(tracerName)

		ctx := extractTraceContext(ss.Context())

		ctx, span := tracer.Start(ctx, info.FullMethod,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("rpc.system", "grpc"),
				attribute.String("rpc.method", info.FullMethod),
				attribute.Bool("rpc.stream", true),
			),
		)
		defer span.End()

		// Propagate request ID
		if reqID := RequestIDFromContext(ctx); reqID != "" {
			span.SetAttributes(attribute.String("request.id", reqID))
		}

		wrapped := &tracedStream{ServerStream: ss, ctx: ctx}
		err := handler(srv, wrapped)
		if err != nil {
			st, _ := status.FromError(err)
			span.SetStatus(codes.Error, st.Message())
			span.SetAttributes(
				attribute.String("rpc.grpc.status_code", st.Code().String()),
			)
			span.RecordError(err)
		} else {
			span.SetStatus(codes.Ok, "")
			span.SetAttributes(
				attribute.String("rpc.grpc.status_code", grpcCodes.OK.String()),
			)
		}

		return err
	}
}

// tracedStream wraps grpc.ServerStream with a traced context.
type tracedStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *tracedStream) Context() context.Context {
	return s.ctx
}

// extractTraceContext extracts W3C trace context from gRPC metadata
// and injects it into the Go context for OpenTelemetry propagation.
func extractTraceContext(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}
	propagator := otel.GetTextMapPropagator()
	return propagator.Extract(ctx, &metadataCarrier{md: md})
}

// metadataCarrier adapts gRPC metadata.MD for OpenTelemetry TextMapPropagator.
type metadataCarrier struct {
	md metadata.MD
}

func (c *metadataCarrier) Get(key string) string {
	vals := c.md.Get(key)
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}

func (c *metadataCarrier) Set(key, value string) {
	c.md.Set(key, value)
}

func (c *metadataCarrier) Keys() []string {
	keys := make([]string, 0, len(c.md))
	for k := range c.md {
		keys = append(keys, k)
	}
	return keys
}

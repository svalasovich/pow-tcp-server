package interceptors

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/svalasovich/golang-template/internal/log"
)

const requestIDHeader = "request-id"

func UnaryServerRequestID() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		if requestIDs := md.Get(requestIDHeader); len(requestIDs) > 0 {
			ctx = context.WithValue(ctx, log.RequestIDContextKey, requestIDs[0])
		} else {
			ctx = context.WithValue(ctx, log.RequestIDContextKey, newRequestID())
		}

		return handler(ctx, req)
	}
}

func UnaryClientRequestID() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		requestID, ok := ctx.Value(log.RequestIDContextKey).(string)
		if !ok || requestID == "" {
			requestID = newRequestID()
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(make(map[string]string))
		}

		md.Set(requestIDHeader, requestID)
		ctx = metadata.NewIncomingContext(ctx, md)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func newRequestID() string {
	return uuid.NewString()
}

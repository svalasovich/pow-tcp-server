package interceptors

import (
	"context"
	"log/slog"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"

	"github.com/svalasovich/golang-template/internal/log"
)

type grpcLogger struct {
	logger *log.Logger
}

func (g *grpcLogger) Log(ctx context.Context, level logging.Level, msg string, fields ...any) {
	g.logger.Log(ctx, slog.Level(level), msg, fields...)
}

func UnaryServerLogging(logger *log.Logger) grpc.UnaryServerInterceptor {
	return logging.UnaryServerInterceptor(&grpcLogger{logger})
}

func UnaryClientLogging(logger *log.Logger) grpc.UnaryClientInterceptor {
	return logging.UnaryClientInterceptor(&grpcLogger{logger})
}

func StreamServerLogging(logger *log.Logger) grpc.StreamServerInterceptor {
	return logging.StreamServerInterceptor(&grpcLogger{logger})
}

func StreamClientLogging(logger *log.Logger) grpc.StreamClientInterceptor {
	return logging.StreamClientInterceptor(&grpcLogger{logger})
}

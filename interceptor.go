package ginney

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"io"
	"time"
)

func CorrelationIdUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		_, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.InvalidArgument, "%s is missing from metadata", CorrelationIdHeaderKey)
		}
		return handler(ctx, req)
	}
}

func LogWithCorrelationIdUnaryServerInterceptor(out io.Writer, ignoreList []string) grpc.UnaryServerInterceptor {
	mustIgnoreLogging := func(apiName string) bool {
		if ignoreList == nil {
			return false
		}

		for _, ignoreItem := range ignoreList {
			if ignoreItem == apiName {
				return true
			}
		}
		return false
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		res, handlerErr := handler(ctx, req)
		if mustIgnoreLogging(info.FullMethod) {
			return res, handlerErr
		}

		// Logging function
		go func() {
			// ip
			ip := "-"
			p, ok := peer.FromContext(ctx)
			if ok {
				ip = p.Addr.String()
			}

			// status code
			statusCode := "-"
			st, ok := status.FromError(handlerErr)
			if ok {
				statusCode = st.Code().String()
			}

			// correlation id
			correlationId := "-"
			if meta, ok := metadata.FromIncomingContext(ctx); ok {
				if correlationIds := meta.Get(CorrelationIdHeaderKey); len(correlationIds) > 0 {
					correlationId = correlationIds[0]
				}
			}

			// api name
			apiName := info.FullMethod

			// end
			end := time.Now()

			// latency
			latency := end.Sub(start)

			_, _ = fmt.Fprint(out, formatLog(
				end,
				correlationId,
				statusCode,
				latency,
				ip,
				apiName,
				grpcRequestBodyToString(req),
			),
			)
		}()

		return res, handlerErr
	}
}

package ginney

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"testing"
	"time"
)

func TestCorrelationIdUnaryServerInterceptor(t *testing.T) {
	req := map[string]interface{}{"id": "1"}
	info := grpc.UnaryServerInfo{}

	t.Run("Happy", func(t *testing.T) {
		incomingCtx := metadata.NewIncomingContext(context.TODO(), metadata.Pairs(CorrelationIdHeaderKey, "random-uuid"))
		interceptor := CorrelationIdUnaryServerInterceptor()
		iface, err := interceptor(incomingCtx, req, &info, func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		})
		assert.NoError(t, err)
		assert.Empty(t, iface)
	})

	t.Run("Error, correlation id is not found in the metadata", func(t *testing.T) {
		incomingCtx := context.TODO()
		interceptor := CorrelationIdUnaryServerInterceptor()
		_, err := interceptor(incomingCtx, req, &info, func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		})
		assert.Error(t, err)
	})
}

func TestLogWithCorrelationIdUnaryServerInterceptor(t *testing.T) {
	req := map[string]interface{}{"id": "1"}
	info := grpc.UnaryServerInfo{FullMethod: "randomMethod"}

	t.Run("Happy", func(t *testing.T) {
		buffer := new(bytes.Buffer)

		incomingCtx := metadata.NewIncomingContext(context.TODO(), metadata.Pairs(CorrelationIdHeaderKey, "random-uuid"))
		interceptor := LogWithCorrelationIdUnaryServerInterceptor(buffer, nil)
		_, err := interceptor(incomingCtx, req, &info, func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		})
		assert.NoError(t, err)

		// Wait for go routine in log
		time.Sleep(1 * time.Second)

		correlationId, statusCode, apiName, payload := extractLogMessage(buffer.String())
		assert.Equal(t, codes.OK.String(), statusCode)
		assert.Equal(t, "randomMethod", apiName)
		assert.Equal(t, `{"id":"1"}`+"\n", payload)
		assert.Equal(t, "random-uuid", correlationId)
	})

	t.Run("Happy, no log when api name is in ignore list", func(t *testing.T) {
		buffer := new(bytes.Buffer)

		incomingCtx := metadata.NewIncomingContext(context.TODO(), metadata.Pairs(CorrelationIdHeaderKey, "random-uuid"))
		interceptor := LogWithCorrelationIdUnaryServerInterceptor(buffer, []string{"randomMethod"})
		_, err := interceptor(incomingCtx, req, &info, func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, nil
		})
		assert.NoError(t, err)

		// Wait for go routine in log
		time.Sleep(1 * time.Second)

		assert.Equal(t, "", buffer.String())
	})

	t.Run("Happy, handler returns error", func(t *testing.T) {
		buffer := new(bytes.Buffer)

		incomingCtx := metadata.NewIncomingContext(context.TODO(), metadata.Pairs(CorrelationIdHeaderKey, "random-uuid"))
		interceptor := LogWithCorrelationIdUnaryServerInterceptor(buffer, nil)
		_, err := interceptor(incomingCtx, req, &info, func(ctx context.Context, req interface{}) (interface{}, error) {
			return nil, status.Error(codes.InvalidArgument, "")
		})
		assert.Error(t, err)

		// Wait for go routine in log
		time.Sleep(1 * time.Second)

		correlationId, statusCode, apiName, payload := extractLogMessage(buffer.String())
		assert.Equal(t, codes.InvalidArgument.String(), statusCode)
		assert.Equal(t, "randomMethod", apiName)
		assert.Equal(t, `{"id":"1"}`+"\n", payload)
		assert.Equal(t, "random-uuid", correlationId)
	})
}

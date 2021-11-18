package ginney

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFromGinContextToContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Happy", func(t *testing.T) {
		response := httptest.NewRecorder()

		ginContext, _ := gin.CreateTestContext(response)
		ginContext.Request, _ = http.NewRequest("POST", "/fake-uri", nil)
		ginContext.Request.Header.Set(CorrelationIdHeaderKey, "random-uuid")
		context := FromGinContextToContext(ginContext)

		assert.NotEmpty(t, context.Value(GinContextKey))
	})
}

func TestFromContextToGinContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Happy", func(t *testing.T) {
		response := httptest.NewRecorder()

		ginContext, _ := gin.CreateTestContext(response)
		ginContext.Request, _ = http.NewRequest("POST", "/fake-uri", nil)
		ginContext.Request.Header.Set(CorrelationIdHeaderKey, "random-uuid")

		context := FromGinContextToContext(ginContext)
		gc, err := FromContextToGinContext(context)

		assert.NoError(t, err)
		assert.NotEmpty(t, gc)
		assert.Equal(t, "random-uuid", gc.Request.Header.Get(CorrelationIdHeaderKey))
	})

	t.Run("Error - context doesn't hold GinContextKey", func(t *testing.T) {
		response := httptest.NewRecorder()

		ginContext, _ := gin.CreateTestContext(response)
		ginContext.Request, _ = http.NewRequest("POST", "/fake-uri", nil)
		ginContext.Request.Header.Set(CorrelationIdHeaderKey, "random-uuid")

		context := ginContext.Request.Context()
		gc, err := FromContextToGinContext(context)

		assert.Error(t, err)
		assert.Nil(t, gc)
	})

	t.Run("Error - context hold GinContextKey but the type is not *gin.Context", func(t *testing.T) {
		response := httptest.NewRecorder()

		ginContext, _ := gin.CreateTestContext(response)
		ginContext.Request, _ = http.NewRequest("POST", "/fake-uri", nil)
		ginContext.Request.Header.Set(CorrelationIdHeaderKey, "random-uuid")
		ctx := context.WithValue(ginContext.Request.Context(), GinContextKey, 1234)
		ginContext.Request = ginContext.Request.WithContext(ctx)

		context := ginContext.Request.Context()
		gc, err := FromContextToGinContext(context)

		assert.Error(t, err)
		assert.Nil(t, gc)
	})
}

func TestNewGrpcOutgoingContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Happy", func(t *testing.T) {
		response := httptest.NewRecorder()

		ginContext, _ := gin.CreateTestContext(response)
		ginContext.Request, _ = http.NewRequest("POST", "/fake-uri", nil)
		ginContext.Request.Header.Set(CorrelationIdHeaderKey, "random-uuid")
		ctx := FromGinContextToContext(ginContext)
		grpcOutCtx := FromContextToGrpcOutgoingContext(ctx)

		meta, _ := metadata.FromOutgoingContext(grpcOutCtx)
		assert.Equal(t, "random-uuid", meta.Get(CorrelationIdHeaderKey)[0])
	})

	t.Run("Happy, no correlation id is set", func(t *testing.T) {
		response := httptest.NewRecorder()

		ginContext, _ := gin.CreateTestContext(response)
		ginContext.Request, _ = http.NewRequest("POST", "/fake-uri", nil)
		ctx := FromGinContextToContext(ginContext)
		grpcOutCtx := FromContextToGrpcOutgoingContext(ctx)

		meta, _ := metadata.FromOutgoingContext(grpcOutCtx)
		assert.Len(t, meta.Get(CorrelationIdHeaderKey), 0)
	})
}

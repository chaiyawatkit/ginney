package ginney

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

func FromContextToGinContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(GinContextKey)
	if ginContext == nil {
		return nil, errors.New("could not retrieve gin.Context")
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		return nil, errors.New("wrong type is in GinContextKey")
	}

	return gc, nil
}

func FromGinContextToContext(c *gin.Context) context.Context {
	ctx := context.WithValue(c.Request.Context(), GinContextKey, c)
	return ctx
}

func FromContextToGrpcOutgoingContext(ctx context.Context) context.Context {
	ginContext, err := FromContextToGinContext(ctx)
	if err != nil {
		return ctx
	}

	correlationId := ginContext.GetHeader(CorrelationIdHeaderKey)
	if correlationId == "" {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, CorrelationIdHeaderKey, correlationId)
}

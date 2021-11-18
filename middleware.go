package ginney

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/nu7hatch/gouuid"
)

func CompositeCorrelationIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationId := c.Request.Header.Get(CorrelationIdHeaderKey)

		if strings.TrimSpace(correlationId) == "" {
			id, _ := uuid.NewV4()
			correlationId = id.String()

			c.Request.Header.Set(CorrelationIdHeaderKey, correlationId)

		}

		c.Writer.Header().Set(CorrelationIdHeaderKey, correlationId)

		c.Next()
	}
}

func MicroServiceCorrelationIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationId := c.Request.Header.Get(CorrelationIdHeaderKey)

		if strings.TrimSpace(correlationId) == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				NewErrorResponse(fmt.Sprintf("%s is missing", CorrelationIdHeaderKey)))
		}

		c.Next()
	}
}

func FromGinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := FromGinContextToContext(c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func LogWithCorrelationIdMiddleware(out io.Writer, notLogged []string) gin.HandlerFunc {
	var skip map[string]struct{}

	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notLogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		if _, ok := skip[path]; !ok {
			correlationId := c.Request.Header.Get(CorrelationIdHeaderKey)

			end := time.Now()
			latency := end.Sub(start)

			clientIP := c.ClientIP()
			method := c.Request.Method

			statusCode := fmt.Sprintf("%d", c.Writer.Status())

			if raw != "" {
				path = path + "?" + raw
			}

			apiName := fmt.Sprintf("%-7s %s", method, path)

			_, _ = fmt.Fprint(out, formatLog(
				end,
				correlationId,
				statusCode,
				latency,
				clientIP,
				apiName,
				httpRequestBodyToString(c.Request.Body),
			),
			)
		}
	}
}

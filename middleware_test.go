package ginney

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
)

func TestLogWithCorrelationIdMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Happy - GET method no query param, no correlation id, enabled logging path", func(t *testing.T) {
		buffer := new(bytes.Buffer)

		// prepare route and register a middleware
		router := gin.New()
		router.Use(LogWithCorrelationIdMiddleware(buffer, []string{}))
		router.GET("/random", func(c *gin.Context) {
			c.AbortWithStatus(http.StatusOK)
		})

		// send request to the route
		_ = performRequest(router, "GET", "/random", nil)
		correlationId, statusCode, path, payload := extractLogMessage(buffer.String())

		assert.Equal(t, strconv.Itoa(http.StatusOK), statusCode)
		assert.Equal(t, fmt.Sprintf("%-7s %s", http.MethodGet, "/random"), path)
		assert.Equal(t, "{}\n", payload)
		assert.Empty(t, correlationId)
	})

	t.Run("Happy - GET method, disabled logging path", func(t *testing.T) {
		buffer := new(bytes.Buffer)

		// prepare route and register a middleware
		router := gin.New()
		router.Use(LogWithCorrelationIdMiddleware(buffer, []string{"/health"}))
		router.GET("/health", func(c *gin.Context) {
			c.AbortWithStatus(http.StatusOK)
		})

		// send request to the route
		_ = performRequest(router, "GET", "/health", nil)

		assert.Empty(t, buffer.String())
	})

	t.Run("Happy - GET method with query params, correlation id, enabled logging path", func(t *testing.T) {
		buffer := new(bytes.Buffer)

		// prepare route and register a middleware
		router := gin.New()
		router.Use(LogWithCorrelationIdMiddleware(buffer, []string{}))
		router.GET("/random", func(c *gin.Context) {
			c.AbortWithStatus(http.StatusOK)
		})

		// send request to the route
		_ = performRequest(router, "GET", "/random?queryParam=123&queryParam=hello%20world", nil,
			header{
				Key:   CorrelationIdHeaderKey,
				Value: "random-uuid",
			},
		)

		correlationId, statusCode, path, payload := extractLogMessage(buffer.String())

		assert.Equal(t, strconv.Itoa(http.StatusOK), statusCode)
		assert.Equal(t, fmt.Sprintf("%-7s %s", http.MethodGet, "/random?queryParam=123&queryParam=hello%20world"), path)
		assert.Equal(t, "{}\n", payload)
		assert.Equal(t, "random-uuid", correlationId)
	})

	t.Run("Happy - GET method with query param, no correlation id, enabled logging path", func(t *testing.T) {
		buffer := new(bytes.Buffer)

		// prepare route and register a middleware
		router := gin.New()
		router.Use(LogWithCorrelationIdMiddleware(buffer, []string{"/health"}))
		router.GET("/random", func(c *gin.Context) {
			c.AbortWithStatus(http.StatusOK)
		})

		// send request to the route
		_ = performRequest(router, "GET", "/random?queryParam=123&queryParam=hello%20world", nil,
			header{
				Key:   CorrelationIdHeaderKey,
				Value: "random-uuid",
			},
		)

		correlationId, statusCode, path, payload := extractLogMessage(buffer.String())

		assert.Equal(t, strconv.Itoa(http.StatusOK), statusCode)
		assert.Equal(t, fmt.Sprintf("%-7s %s", http.MethodGet, "/random?queryParam=123&queryParam=hello%20world"), path)
		assert.Equal(t, "{}\n", payload)
		assert.Equal(t, "random-uuid", correlationId)
	})

	t.Run("Happy - POST method with Body, correlation id, enabled logging path", func(t *testing.T) {
		buffer := new(bytes.Buffer)

		// prepare route and register a middleware
		router := gin.New()
		router.Use(LogWithCorrelationIdMiddleware(buffer, []string{"/health"}))
		router.POST("/random", func(c *gin.Context) {
			c.AbortWithStatus(http.StatusOK)
		})

		requestBody, _ := json.Marshal(randomJson{
			Example: "hello",
		})

		// send request to the route
		_ = performRequest(router, "POST", "/random",
			bytes.NewBuffer(requestBody),
			header{
				Key:   CorrelationIdHeaderKey,
				Value: "random-uuid",
			},
		)

		correlationId, statusCode, path, payload := extractLogMessage(buffer.String())

		assert.Equal(t, strconv.Itoa(http.StatusOK), statusCode)
		assert.Equal(t, fmt.Sprintf("%-7s %s", http.MethodPost, "/random"), path)
		assert.Equal(t, `{"example":"hello"}`+"\n", payload)
		assert.Equal(t, "random-uuid", correlationId)
	})

	t.Run("Happy - POST method without Body, correlation id, enabled logging path", func(t *testing.T) {
		buffer := new(bytes.Buffer)

		// prepare route and register a middleware
		router := gin.New()
		router.Use(LogWithCorrelationIdMiddleware(buffer, []string{"/health"}))
		router.POST("/random", func(c *gin.Context) {
			c.AbortWithStatus(http.StatusOK)
		})

		// send request to the route
		_ = performRequest(router, "POST", "/random", nil,
			header{
				Key:   CorrelationIdHeaderKey,
				Value: "random-uuid",
			},
		)

		correlationId, statusCode, path, payload := extractLogMessage(buffer.String())

		assert.Equal(t, strconv.Itoa(http.StatusOK), statusCode)
		assert.Equal(t, fmt.Sprintf("%-7s %s", http.MethodPost, "/random"), path)
		assert.Equal(t, "{}\n", payload)
		assert.Equal(t, "random-uuid", correlationId)
	})

	t.Run("Happy - POST method with Body, no correlation id, enabled logging path", func(t *testing.T) {
		buffer := new(bytes.Buffer)

		// prepare route and register a middleware
		router := gin.New()
		router.Use(LogWithCorrelationIdMiddleware(buffer, []string{"/health"}))
		router.POST("/random", func(c *gin.Context) {
			c.AbortWithStatus(http.StatusOK)
		})

		requestBody, _ := json.Marshal(randomJson{
			Example: "hello",
		})

		// send request to the route
		_ = performRequest(router, "POST", "/random",
			bytes.NewBuffer(requestBody),
		)

		correlationId, statusCode, path, payload := extractLogMessage(buffer.String())

		assert.Equal(t, strconv.Itoa(http.StatusOK), statusCode)
		assert.Equal(t, fmt.Sprintf("%-7s %s", http.MethodPost, "/random"), path)
		assert.Equal(t, `{"example":"hello"}`+"\n", payload)
		assert.Empty(t, correlationId)
	})

	t.Run("Happy - POST method without Body, no correlation id, enabled logging path", func(t *testing.T) {
		buffer := new(bytes.Buffer)

		// prepare route and register a middleware
		router := gin.New()
		router.Use(LogWithCorrelationIdMiddleware(buffer, []string{"/health"}))
		router.POST("/random", func(c *gin.Context) {
			c.AbortWithStatus(http.StatusOK)
		})

		// send request to the route
		_ = performRequest(router, "POST", "/random", nil)

		correlationId, statusCode, path, payload := extractLogMessage(buffer.String())

		assert.Equal(t, strconv.Itoa(http.StatusOK), statusCode)
		assert.Equal(t, fmt.Sprintf("%-7s %s", http.MethodPost, "/random"), path)
		assert.Equal(t, "{}\n", payload)
		assert.Empty(t, correlationId)
	})

	t.Run("Happy - Put method with Body, correlation id, enabled logging path", func(t *testing.T) {
		buffer := new(bytes.Buffer)

		// prepare route and register a middleware
		router := gin.New()
		router.Use(LogWithCorrelationIdMiddleware(buffer, []string{"/health"}))
		router.PUT("/random", func(c *gin.Context) {
			c.AbortWithStatus(http.StatusOK)
		})

		requestBody, _ := json.Marshal(randomJson{
			Example: "hello",
		})

		// send request to the route
		_ = performRequest(router, "PUT", "/random",
			bytes.NewBuffer(requestBody),
			header{
				Key:   CorrelationIdHeaderKey,
				Value: "random-uuid",
			},
		)

		correlationId, statusCode, path, payload := extractLogMessage(buffer.String())

		assert.Equal(t, strconv.Itoa(http.StatusOK), statusCode)
		assert.Equal(t, fmt.Sprintf("%-7s %s", http.MethodPut, "/random"), path)
		assert.Equal(t, `{"example":"hello"}`+"\n", payload)
		assert.Equal(t, "random-uuid", correlationId)
	})

	t.Run("Happy - PUT method without Body, correlation id, enabled logging path", func(t *testing.T) {
		buffer := new(bytes.Buffer)

		// prepare route and register a middleware
		router := gin.New()
		router.Use(LogWithCorrelationIdMiddleware(buffer, []string{"/health"}))
		router.PUT("/random", func(c *gin.Context) {
			c.AbortWithStatus(http.StatusOK)
		})

		// send request to the route
		_ = performRequest(router, "PUT", "/random", nil,
			header{
				Key:   CorrelationIdHeaderKey,
				Value: "random-uuid",
			},
		)

		correlationId, statusCode, path, payload := extractLogMessage(buffer.String())

		assert.Equal(t, strconv.Itoa(http.StatusOK), statusCode)
		assert.Equal(t, fmt.Sprintf("%-7s %s", http.MethodPut, "/random"), path)
		assert.Equal(t, "{}\n", payload)
		assert.Equal(t, "random-uuid", correlationId)
	})

	t.Run("Happy - PUT method with Body, no correlation id, enabled logging path", func(t *testing.T) {
		buffer := new(bytes.Buffer)

		// prepare route and register a middleware
		router := gin.New()
		router.Use(LogWithCorrelationIdMiddleware(buffer, []string{"/health"}))
		router.PUT("/random", func(c *gin.Context) {
			c.AbortWithStatus(http.StatusOK)
		})

		requestBody, _ := json.Marshal(randomJson{
			Example: "hello",
		})

		// send request to the route
		_ = performRequest(router, "PUT", "/random",
			bytes.NewBuffer(requestBody),
		)

		correlationId, statusCode, path, payload := extractLogMessage(buffer.String())

		assert.Equal(t, strconv.Itoa(http.StatusOK), statusCode)
		assert.Equal(t, fmt.Sprintf("%-7s %s", http.MethodPut, "/random"), path)
		assert.Equal(t, `{"example":"hello"}`+"\n", payload)
		assert.Empty(t, correlationId)
	})

	t.Run("Happy - PUT method without Body, no correlation id, enabled logging path", func(t *testing.T) {
		buffer := new(bytes.Buffer)

		// prepare route and register a middleware
		router := gin.New()
		router.Use(LogWithCorrelationIdMiddleware(buffer, []string{"/health"}))
		router.PUT("/random", func(c *gin.Context) {
			c.AbortWithStatus(http.StatusOK)
		})

		// send request to the route
		_ = performRequest(router, "PUT", "/random", nil)

		correlationId, statusCode, path, payload := extractLogMessage(buffer.String())

		assert.Equal(t, strconv.Itoa(http.StatusOK), statusCode)
		assert.Equal(t, fmt.Sprintf("%-7s %s", http.MethodPut, "/random"), path)
		assert.Equal(t, "{}\n", payload)
		assert.Empty(t, correlationId)
	})

	t.Run("Happy - Test Body with censored field", func(t *testing.T) {
		buffer := new(bytes.Buffer)

		// prepare route and register a middleware
		router := gin.New()
		router.Use(LogWithCorrelationIdMiddleware(buffer, []string{"/health"}))
		router.POST("/random", func(c *gin.Context) {
			c.AbortWithStatus(http.StatusOK)
		})

		requestBody, _ := json.Marshal(randomJson{
			Example:           "hello",
			ExamplePassword:   "V3l0S3cr3t",
			ExamplePrivateKey: "privkey.hahaha",
			ExampleSecretKey:  "skey.hahaha",
		})

		// send request to the route
		_ = performRequest(router, "POST", "/random",
			bytes.NewBuffer(requestBody),
		)

		_, _, _, payload := extractLogMessage(buffer.String())

		assert.Equal(t, `{"example":"hello","examplePassword":"[HIDDEN_FIELD]","examplePrivateKey":"[HIDDEN_FIELD]","exampleSecretKey":"[HIDDEN_FIELD]"}`+"\n", payload)
	})

}

func TestFromGinContextToContextMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Happy - with middleware correctly setup", func(t *testing.T) {
		// prepare route and register a middleware
		router := gin.New()
		router.Use(FromGinContextToContextMiddleware())
		router.GET("/random", func(c *gin.Context) {
			assert.NotEmpty(t, c.Request.Context().Value(GinContextKey))
			c.AbortWithStatus(http.StatusOK)
		})

		// send request to the route
		_ = performRequest(router, "GET", "/random", nil)
	})

	t.Run("Happy - without middleware", func(t *testing.T) {
		// prepare route and register a middleware
		router := gin.New()
		router.GET("/random", func(c *gin.Context) {
			assert.Empty(t, c.Request.Context().Value(GinContextKey))
			c.AbortWithStatus(http.StatusOK)
		})

		// send request to the route
		_ = performRequest(router, "GET", "/random", nil)
	})
}

func TestCompositeCorrelationIdMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Happy - correlation id is provided", func(t *testing.T) {
		// prepare route and register a middleware
		router := gin.New()
		router.Use(CompositeCorrelationIdMiddleware())
		router.POST("/random", func(c *gin.Context) {
			assert.NotEmpty(t, c.Request.Header.Get(CorrelationIdHeaderKey))
			assert.NotEmpty(t, c.Writer.Header().Get(CorrelationIdHeaderKey))
			assert.Equal(t, "random-uuid", c.Request.Header.Get(CorrelationIdHeaderKey))
			assert.Equal(t, "random-uuid", c.Writer.Header().Get(CorrelationIdHeaderKey))

			c.AbortWithStatus(http.StatusOK)
		})

		// send request to the route
		_ = performRequest(router, "POST", "/random", nil,
			header{
				Key:   CorrelationIdHeaderKey,
				Value: "random-uuid",
			},
		)
	})

	t.Run("Happy - correlation id isn't provided", func(t *testing.T) {
		// prepare route and register a middleware
		router := gin.New()
		router.Use(CompositeCorrelationIdMiddleware())
		router.POST("/random", func(c *gin.Context) {
			assert.NotEmpty(t, c.Request.Header.Get(CorrelationIdHeaderKey))
			assert.NotEmpty(t, c.Writer.Header().Get(CorrelationIdHeaderKey))

			c.AbortWithStatus(http.StatusOK)
		})

		// send request to the route
		_ = performRequest(router, "POST", "/random", nil)
	})
}

func TestMicroServiceCorrelationIdMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Happy - correlation id is provided", func(t *testing.T) {
		// prepare route and register a middleware
		router := gin.New()
		router.Use(MicroServiceCorrelationIdMiddleware())
		router.POST("/random", func(c *gin.Context) {
			assert.NotEmpty(t, c.Request.Header.Get(CorrelationIdHeaderKey))
			assert.Equal(t, "random-uuid", c.Request.Header.Get(CorrelationIdHeaderKey))

			c.AbortWithStatus(http.StatusOK)
		})

		// send request to the route
		_ = performRequest(router, "POST", "/random", nil,
			header{
				Key:   CorrelationIdHeaderKey,
				Value: "random-uuid",
			},
		)
	})

	t.Run("Error - correlation id isn't provided", func(t *testing.T) {
		// prepare route and register a middleware
		router := gin.New()
		router.Use(MicroServiceCorrelationIdMiddleware())
		router.POST("/random", func(c *gin.Context) {
			c.AbortWithStatus(http.StatusOK)
		})

		// send request to the route
		resp := performRequest(router, "POST", "/random", nil)

		var respBody map[string]interface{}
		err := json.Unmarshal(resp.Body.Bytes(), &respBody)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.Equal(t, StatusFail, respBody["status"])
		assert.Equal(t, fmt.Sprintf("%s is missing", CorrelationIdHeaderKey), respBody["message"])
	})
}

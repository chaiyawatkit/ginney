package ginney

import (
	"github.com/gin-gonic/gin"
	"github.com/jarcoal/httpmock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

type header struct {
	Key   string
	Value string
}

type randomJson struct {
	Example           string `json:"example"`
	ExamplePassword   string `json:"examplePassword,omitempty"`
	ExamplePrivateKey string `json:"examplePrivateKey,omitempty"`
	ExampleSecretKey  string `json:"exampleSecretKey,omitempty"`
}

func createGinContextWithCorrelationId(method string, url string, correlationId string) *gin.Context {
	gin.SetMode(gin.TestMode)

	ginContext, _ := gin.CreateTestContext(httptest.NewRecorder())
	ginContext.Request, _ = http.NewRequest(method, url, nil)
	ginContext.Request.Header.Set(CorrelationIdHeaderKey, correlationId)

	return ginContext
}

func initHttpMock(httpMethod string, url string, statusCode int, fixture string) {
	fakeRespond := httpmock.NewStringResponder(statusCode, fixture)
	httpmock.Activate()
	httpmock.RegisterResponder(httpMethod, url, fakeRespond)
}

func performRequest(r http.Handler, method, path string, body io.Reader, headers ...header) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func extractLogMessage(logMsg string) (correlationId string, statusCode string, apiName string, payload string) {
	cols := strings.Split(logMsg, "|")

	for index, elem := range cols {
		cols[index] = strings.Trim(elem, " ")
	}

	return cols[1], cols[2], cols[5], cols[6]
}

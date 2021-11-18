package ginney

import (
	"context"
	"io"
	"net/http"
)

func send(ctx context.Context, method string, url string, contentType string, body io.Reader) (*http.Response, error) {
	ginContext, _ := FromContextToGinContext(ctx)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// setting headers
	if ginContext != nil {
		req.Header.Set(CorrelationIdHeaderKey, ginContext.GetHeader(CorrelationIdHeaderKey))
	}
	if contentType != "" {
		req.Header.Set(ContentTypeHeaderKey, contentType)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func Post(ctx context.Context, url string, contentType string, body io.Reader) (*http.Response, error) {
	return send(ctx, http.MethodPost, url, contentType, body)
}

func Get(ctx context.Context, url string) (*http.Response, error) {
	return send(ctx, http.MethodGet, url, "", nil)
}

func Put(ctx context.Context, url string, contentType string, body io.Reader) (*http.Response, error) {
	return send(ctx, http.MethodPut, url, contentType, body)
}

func Delete(ctx context.Context, url string, contentType string, body io.Reader) (*http.Response, error) {
	return send(ctx, http.MethodDelete, url, contentType, body)
}

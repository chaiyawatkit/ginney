package ginney

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

func TestGet(t *testing.T) {
	t.Run("Happy", func(t *testing.T) {
		initHttpMock(http.MethodGet, "https://www.fcuk.com", http.StatusOK, `{"ping": "pong"}`)

		gc := createGinContextWithCorrelationId(http.MethodGet, "/chaiyawatkit", "random-uuid")
		ctx := context.WithValue(context.TODO(), GinContextKey, gc)

		resp, err := Get(ctx, "https://www.fcuk.com")
		assert.NoError(t, err)

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		var respBody map[string]interface{}
		err = json.Unmarshal(respBodyBytes, &respBody)
		assert.NoError(t, err)

		assert.NotEmpty(t, resp)
		assert.Equal(t, "pong", respBody["ping"])
	})

	t.Run("Happy - GinContextKey contains wrong data type", func(t *testing.T) {
		initHttpMock(http.MethodGet, "https://www.fcuk.com", http.StatusOK, `{"ping": "pong"}`)

		ctx := context.WithValue(context.TODO(), GinContextKey, 1234)

		resp, err := Get(ctx, "https://www.fcuk.com")
		assert.NoError(t, err)

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		var respBody map[string]interface{}
		err = json.Unmarshal(respBodyBytes, &respBody)
		assert.NoError(t, err)

		assert.NoError(t, err)
		assert.Equal(t, "pong", respBody["ping"])
	})

	t.Run("Happy - context doesn't contain GinContextKey at all", func(t *testing.T) {
		initHttpMock(http.MethodGet, "https://www.fcuk.com", http.StatusOK, `{"ping": "pong"}`)

		ctx := context.TODO()

		resp, err := Get(ctx, "https://www.fcuk.com")
		assert.NoError(t, err)

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		var respBody map[string]interface{}
		err = json.Unmarshal(respBodyBytes, &respBody)
		assert.NoError(t, err)
		assert.Equal(t, "pong", respBody["ping"])
	})

	t.Run("Error - bad request", func(t *testing.T) {
		initHttpMock(http.MethodGet, "https://www.fcuk.com", http.StatusBadRequest, `{"ping": "pong"}`)

		gc := createGinContextWithCorrelationId(http.MethodGet, "/chaiyawatkit", "random-uuid")
		ctx := context.WithValue(context.TODO(), GinContextKey, gc)

		resp, err := Get(ctx, "https://www.fcuk.com")

		assert.NoError(t, err)
		assert.NotEmpty(t, resp)
		assert.Equal(t, strconv.Itoa(http.StatusBadRequest), resp.Status)
	})

	t.Run("Happy - context GinContextKey but there's no correlation id in the header", func(t *testing.T) {
		initHttpMock(http.MethodGet, "https://www.fcuk.com", http.StatusOK, `{"ping": "pong"}`)

		gc := createGinContextWithCorrelationId(http.MethodGet, "/chaiyawatkit", "")
		ctx := context.WithValue(context.TODO(), GinContextKey, gc)

		resp, err := Get(ctx, "https://www.fcuk.com")
		assert.NoError(t, err)

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		var respBody map[string]interface{}
		err = json.Unmarshal(respBodyBytes, &respBody)
		assert.NoError(t, err)
		assert.Equal(t, "pong", respBody["ping"])
	})
}

func TestPost(t *testing.T) {
	t.Run("Happy - without body", func(t *testing.T) {
		initHttpMock(http.MethodPost, "https://www.fcuk.com", http.StatusOK, `{"ping": "pong"}`)

		gc := createGinContextWithCorrelationId(http.MethodPost, "/chaiyawatkit", "random-uuid")
		ctx := context.WithValue(context.TODO(), GinContextKey, gc)

		resp, err := Post(ctx, "https://www.fcuk.com", "application/json", nil)
		assert.NoError(t, err)

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		var respBody map[string]interface{}
		err = json.Unmarshal(respBodyBytes, &respBody)
		assert.NoError(t, err)

		assert.NotEmpty(t, resp)
		assert.Equal(t, "pong", respBody["ping"])
	})

	t.Run("Happy - with body", func(t *testing.T) {
		initHttpMock(http.MethodPost, "https://www.fcuk.com", http.StatusOK, `{"ping": "pong"}`)

		gc := createGinContextWithCorrelationId(http.MethodPost, "/chaiyawatkit", "random-uuid")
		ctx := context.WithValue(context.TODO(), GinContextKey, gc)

		requestBody, _ := json.Marshal(randomJson{
			Example: "hello",
		})

		resp, err := Post(ctx, "https://www.fcuk.com", "application/json", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		var respBody map[string]interface{}
		err = json.Unmarshal(respBodyBytes, &respBody)
		assert.NoError(t, err)

		assert.NotEmpty(t, resp)
		assert.Equal(t, "pong", respBody["ping"])
	})

	t.Run("Happy - GinContextKey contains wrong data type", func(t *testing.T) {
		initHttpMock(http.MethodPost, "https://www.fcuk.com", http.StatusOK, `{"ping": "pong"}`)

		ctx := context.WithValue(context.TODO(), GinContextKey, 1234)

		resp, err := Post(ctx, "https://www.fcuk.com", "application/json", nil)
		assert.NoError(t, err)

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		var respBody map[string]interface{}
		err = json.Unmarshal(respBodyBytes, &respBody)
		assert.NoError(t, err)

		assert.NotEmpty(t, resp)
		assert.Equal(t, "pong", respBody["ping"])
	})

	t.Run("Happy - context doesn't contain GinContextKey at all", func(t *testing.T) {
		initHttpMock(http.MethodGet, "https://www.fcuk.com", http.StatusOK, `{"ping": "pong"}`)

		ctx := context.TODO()

		resp, err := Post(ctx, "https://www.fcuk.com", "application/json", nil)
		assert.NoError(t, err)

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		var respBody map[string]interface{}
		err = json.Unmarshal(respBodyBytes, &respBody)
		assert.NoError(t, err)

		assert.NotEmpty(t, resp)
		assert.Equal(t, "pong", respBody["ping"])
	})

	t.Run("Error - bad request", func(t *testing.T) {
		initHttpMock(http.MethodPost, "https://www.fcuk.com", http.StatusBadRequest, `{"ping": "pong"}`)

		gc := createGinContextWithCorrelationId(http.MethodPost, "/chaiyawatkit", "random-uuid")
		ctx := context.WithValue(context.TODO(), GinContextKey, gc)

		resp, err := Post(ctx, "https://www.fcuk.com", "application/json", nil)

		assert.NoError(t, err)
		assert.NotEmpty(t, resp)
		assert.Equal(t, strconv.Itoa(http.StatusBadRequest), resp.Status)
	})

	t.Run("Happy - context GinContextKey but there's no correlation id in the header", func(t *testing.T) {
		initHttpMock(http.MethodPost, "https://www.fcuk.com", http.StatusOK, `{"ping": "pong"}`)

		gc := createGinContextWithCorrelationId(http.MethodPost, "/chaiyawatkit", "")
		ctx := context.WithValue(context.TODO(), GinContextKey, gc)

		resp, err := Post(ctx, "https://www.fcuk.com", "application/json", nil)
		assert.NoError(t, err)

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		var respBody map[string]interface{}
		err = json.Unmarshal(respBodyBytes, &respBody)
		assert.NoError(t, err)
		assert.Equal(t, "pong", respBody["ping"])
	})
}

func TestPut(t *testing.T) {
	t.Run("Happy - without body", func(t *testing.T) {
		initHttpMock(http.MethodPut, "https://www.fcuk.com", http.StatusOK, `{"ping": "pong"}`)

		gc := createGinContextWithCorrelationId(http.MethodPut, "/chaiyawatkit", "random-uuid")
		ctx := context.WithValue(context.TODO(), GinContextKey, gc)

		resp, err := Put(ctx, "https://www.fcuk.com", "application/json", nil)
		assert.NoError(t, err)

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		var respBody map[string]interface{}
		err = json.Unmarshal(respBodyBytes, &respBody)
		assert.NoError(t, err)

		assert.NotEmpty(t, resp)
		assert.Equal(t, "pong", respBody["ping"])
	})

	t.Run("Happy - with body", func(t *testing.T) {
		initHttpMock(http.MethodPut, "https://www.fcuk.com", http.StatusOK, `{"ping": "pong"}`)

		gc := createGinContextWithCorrelationId(http.MethodPut, "/chaiyawatkit", "random-uuid")
		ctx := context.WithValue(context.TODO(), GinContextKey, gc)

		requestBody, _ := json.Marshal(randomJson{
			Example: "hello",
		})

		resp, err := Put(ctx, "https://www.fcuk.com", "application/json", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		var respBody map[string]interface{}
		err = json.Unmarshal(respBodyBytes, &respBody)
		assert.NoError(t, err)

		assert.NotEmpty(t, resp)
		assert.Equal(t, "pong", respBody["ping"])
	})

	t.Run("Happy - GinContextKey contains wrong data type", func(t *testing.T) {
		initHttpMock(http.MethodPut, "https://www.fcuk.com", http.StatusOK, `{"ping": "pong"}`)

		ctx := context.WithValue(context.TODO(), GinContextKey, 1234)

		resp, err := Put(ctx, "https://www.fcuk.com", "application/json", nil)
		assert.NoError(t, err)

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		var respBody map[string]interface{}
		err = json.Unmarshal(respBodyBytes, &respBody)
		assert.NoError(t, err)

		assert.NotEmpty(t, resp)
		assert.Equal(t, "pong", respBody["ping"])
	})

	t.Run("Happy - context doesn't contain GinContextKey at all", func(t *testing.T) {
		initHttpMock(http.MethodPut, "https://www.fcuk.com", http.StatusOK, `{"ping": "pong"}`)

		ctx := context.TODO()

		resp, err := Put(ctx, "https://www.fcuk.com", "application/json", nil)
		assert.NoError(t, err)

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		var respBody map[string]interface{}
		err = json.Unmarshal(respBodyBytes, &respBody)
		assert.NoError(t, err)

		assert.NotEmpty(t, resp)
		assert.Equal(t, "pong", respBody["ping"])
	})

	t.Run("Error - bad request", func(t *testing.T) {
		initHttpMock(http.MethodPut, "https://www.fcuk.com", http.StatusBadRequest, `{"ping": "pong"}`)

		gc := createGinContextWithCorrelationId(http.MethodPut, "/chaiyawatkit", "random-uuid")
		ctx := context.WithValue(context.TODO(), GinContextKey, gc)

		resp, err := Put(ctx, "https://www.fcuk.com", "application/json", nil)

		assert.NoError(t, err)
		assert.NotEmpty(t, resp)
		assert.Equal(t, strconv.Itoa(http.StatusBadRequest), resp.Status)
	})

	t.Run("Happy - context GinContextKey but there's no correlation id in the header", func(t *testing.T) {
		initHttpMock(http.MethodPut, "https://www.fcuk.com", http.StatusOK, `{"ping": "pong"}`)

		gc := createGinContextWithCorrelationId(http.MethodPut, "/chaiyawatkit", "")
		ctx := context.WithValue(context.TODO(), GinContextKey, gc)

		resp, err := Put(ctx, "https://www.fcuk.com", "application/json", nil)
		assert.NoError(t, err)

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		var respBody map[string]interface{}
		err = json.Unmarshal(respBodyBytes, &respBody)
		assert.NoError(t, err)
		assert.Equal(t, "pong", respBody["ping"])
	})
}

func TestDelete(t *testing.T) {
	t.Run("Happy - without body", func(t *testing.T) {
		initHttpMock(http.MethodDelete, "https://www.fcuk.com", http.StatusOK, `{"ping": "pong"}`)

		gc := createGinContextWithCorrelationId(http.MethodDelete, "/chaiyawatkit", "random-uuid")
		ctx := context.WithValue(context.TODO(), GinContextKey, gc)

		resp, err := Delete(ctx, "https://www.fcuk.com", "application/json", nil)
		assert.NoError(t, err)

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		var respBody map[string]interface{}
		err = json.Unmarshal(respBodyBytes, &respBody)
		assert.NoError(t, err)

		assert.NotEmpty(t, resp)
		assert.Equal(t, "pong", respBody["ping"])
	})

	t.Run("Happy - with body", func(t *testing.T) {
		initHttpMock(http.MethodDelete, "https://www.fcuk.com", http.StatusOK, `{"ping": "pong"}`)

		gc := createGinContextWithCorrelationId(http.MethodDelete, "/chaiyawatkit", "random-uuid")
		ctx := context.WithValue(context.TODO(), GinContextKey, gc)

		requestBody, _ := json.Marshal(randomJson{
			Example: "hello",
		})

		resp, err := Delete(ctx, "https://www.fcuk.com", "application/json", bytes.NewBuffer(requestBody))
		assert.NoError(t, err)

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		var respBody map[string]interface{}
		err = json.Unmarshal(respBodyBytes, &respBody)
		assert.NoError(t, err)

		assert.NotEmpty(t, resp)
		assert.Equal(t, "pong", respBody["ping"])
	})

	t.Run("Happy - GinContextKey contains wrong data type", func(t *testing.T) {
		initHttpMock(http.MethodDelete, "https://www.fcuk.com", http.StatusOK, `{"ping": "pong"}`)

		ctx := context.WithValue(context.TODO(), GinContextKey, 1234)

		resp, err := Delete(ctx, "https://www.fcuk.com", "application/json", nil)
		assert.NoError(t, err)

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		var respBody map[string]interface{}
		err = json.Unmarshal(respBodyBytes, &respBody)
		assert.NoError(t, err)

		assert.NotEmpty(t, resp)
		assert.Equal(t, "pong", respBody["ping"])
	})

	t.Run("Happy - context doesn't contain GinContextKey at all", func(t *testing.T) {
		initHttpMock(http.MethodDelete, "https://www.fcuk.com", http.StatusOK, `{"ping": "pong"}`)

		ctx := context.TODO()

		resp, err := Delete(ctx, "https://www.fcuk.com", "application/json", nil)
		assert.NoError(t, err)

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		var respBody map[string]interface{}
		err = json.Unmarshal(respBodyBytes, &respBody)
		assert.NoError(t, err)

		assert.NotEmpty(t, resp)
		assert.Equal(t, "pong", respBody["ping"])
	})

	t.Run("Error - bad request", func(t *testing.T) {
		initHttpMock(http.MethodDelete, "https://www.fcuk.com", http.StatusBadRequest, `{"ping": "pong"}`)

		gc := createGinContextWithCorrelationId(http.MethodDelete, "/chaiyawatkit", "random-uuid")
		ctx := context.WithValue(context.TODO(), GinContextKey, gc)

		resp, err := Delete(ctx, "https://www.fcuk.com", "application/json", nil)

		assert.NoError(t, err)
		assert.NotEmpty(t, resp)
		assert.Equal(t, strconv.Itoa(http.StatusBadRequest), resp.Status)
	})

	t.Run("Happy - context GinContextKey but there's no correlation id in the header", func(t *testing.T) {
		initHttpMock(http.MethodDelete, "https://www.fcuk.com", http.StatusOK, `{"ping": "pong"}`)

		gc := createGinContextWithCorrelationId(http.MethodDelete, "/chaiyawatkit", "")
		ctx := context.WithValue(context.TODO(), GinContextKey, gc)

		resp, err := Delete(ctx, "https://www.fcuk.com", "application/json", nil)
		assert.NoError(t, err)

		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)

		var respBody map[string]interface{}
		err = json.Unmarshal(respBodyBytes, &respBody)
		assert.NoError(t, err)
		assert.Equal(t, "pong", respBody["ping"])
	})
}

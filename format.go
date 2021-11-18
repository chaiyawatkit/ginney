package ginney

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

func formatLog(eventTime time.Time, correlationId string, statusCode string, latency time.Duration, clientIp, apiName, body string) string {
	return fmt.Sprintf("[chaiyawatkit] %v | %5s | %3s | %13v | %15s | %s | %s\n",
		eventTime.Format("2006/01/02 - 15:04:05"),
		correlationId,
		statusCode,
		latency,
		clientIp,
		apiName,
		body,
	)
}

func httpRequestBodyToString(body io.ReadCloser) string {
	if body == nil {
		return ""
	}

	var bodyData map[string]interface{}
	err := json.NewDecoder(body).Decode(&bodyData)
	if err != nil {
		return fmt.Sprintf("%s", body)
	}
	return jsonBodyToString(bodyData)
}

func grpcRequestBodyToString(body interface{}) string {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Sprintf("%s", body)
	}

	var bodyData map[string]interface{}
	err = json.Unmarshal(jsonBytes, &bodyData)
	if err != nil {
		return fmt.Sprintf("%s", body)
	}

	return jsonBodyToString(bodyData)
}

func jsonBodyToString(jsonBody map[string]interface{}) string {
	for key := range jsonBody {
		if shouldKeyCensored(key) {
			jsonBody[key] = CensoredFieldText
		}
	}

	jsonBytes, _ := json.Marshal(jsonBody)
	return string(jsonBytes)
}

func shouldKeyCensored(key string) bool {
	loweredKey := strings.ToLower(key)
	for _, censoredKeyWord := range RequestBodyKeyCensoredList {
		if strings.Contains(loweredKey, censoredKeyWord) {
			return true
		}
	}
	return false
}

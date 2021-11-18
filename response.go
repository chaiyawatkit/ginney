package ginney

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	StatusFail    string = "fail"
	StatusSuccess string = "success"
)

const (
	MessageOk string = "OK"
)

type BaseSuccessResponse struct {
	Status  string      `json:"status" example:"success"`
	Message string      `json:"message" example:"OK"`
	Data    interface{} `json:"data"`
}

func JSONSuccessResponse(c *gin.Context, data interface{}) {
	r := new(BaseSuccessResponse)
	r.Status = StatusSuccess
	r.Message = MessageOk
	r.Data = data
	c.AbortWithStatusJSON(http.StatusOK, *r)
}

func NewSuccessResponse(data interface{}) BaseSuccessResponse {
	r := BaseSuccessResponse{
		Status:  StatusSuccess,
		Message: MessageOk,
		Data:    data,
	}
	return r
}

type ErrorResponse struct {
	Status  string `json:"status" example:"fail"`
	Message string `json:"message" example:"Error message will be show here"`
}

func NewErrorResponse(message string) ErrorResponse {
	errorResponse := ErrorResponse{}
	errorResponse.Status = StatusFail
	errorResponse.Message = message
	return errorResponse
}

type successHumanResponse struct {
	Status       string      `json:"status" example:"fail"`
	Message      string      `json:"message" example:"Error message will be show here"`
	HumanMessage *string     `json:"humanMessage" example:"Error message will be show here"`
	Data         interface{} `json:"data"`
}

type errorHumanResponse struct {
	Status       string  `json:"status" example:"fail"`
	Message      string  `json:"message" example:"Error message will be show here"`
	HumanMessage *string `json:"humanMessage" example:"Error message will be show here"`
}

func JSONErrorHumanResponse(c *gin.Context, errCode int, err error, humanMsg *string) {

	errorResponse := errorHumanResponse{
		Status:  StatusFail,
		Message: err.Error(),
	}
	if humanMsg != nil {
		errorResponse.HumanMessage = humanMsg
	}

	c.AbortWithStatusJSON(errCode, errorResponse)
}

func JSONSuccessHumanResponse(c *gin.Context, data interface{}, humanMsg *string) {

	successResponse := successHumanResponse{
		Status:  StatusSuccess,
		Message: MessageOk,
		Data:    data,
	}
	if humanMsg != nil {
		successResponse.HumanMessage = humanMsg
	}
	c.AbortWithStatusJSON(http.StatusOK, successResponse)
}

type successCodeResponse struct {
	Status  string      `json:"status" example:"fail"`
	Code    string      `json:"code" example:"200"`
	Message string      `json:"message" example:"Error message will be show here"`
	Data    interface{} `json:"data"`
}

type errorCodeResponse struct {
	Status  string `json:"status" example:"fail"`
	Code    string `json:"code" example:"99"`
	Message string `json:"message" example:"Error message will be show here"`
}

func JSONErrorCodeResponse(c *gin.Context, errCode int, err error, code string) {

	errorResponse := errorCodeResponse{
		Status:  StatusFail,
		Message: err.Error(),
		Code:    code,
	}

	c.AbortWithStatusJSON(errCode, errorResponse)
}

func JSONSuccessCodeResponse(c *gin.Context, data interface{}, code string) {

	successResponse := successCodeResponse{
		Status:  StatusSuccess,
		Message: MessageOk,
		Data:    data,
		Code:    code,
	}

	c.AbortWithStatusJSON(http.StatusOK, successResponse)
}

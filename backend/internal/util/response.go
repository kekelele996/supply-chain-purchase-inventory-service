package util

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/supplychain/supplychain-api/internal/constants"
)

// Body 统一三段式响应体。
type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// OK 返回成功响应（code=0, message=ok）。
func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Body{Code: constants.CodeOK, Message: constants.MsgOK, Data: data})
}

// OKWithMessage 返回携带自定义成功文案的响应。
func OKWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Body{Code: constants.CodeOK, Message: message, Data: data})
}

// Created 返回创建成功响应。
func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Body{Code: constants.CodeOK, Message: message, Data: data})
}

// Fail 返回错误响应。
func Fail(c *gin.Context, httpStatus, code int, message string) {
	c.AbortWithStatusJSON(httpStatus, Body{Code: code, Message: message, Data: nil})
}

// FailWithData 返回带数据的错误响应（如字段级校验错误）。
func FailWithData(c *gin.Context, httpStatus, code int, message string, data interface{}) {
	c.AbortWithStatusJSON(httpStatus, Body{Code: code, Message: message, Data: data})
}

// Error 将任意 error 转为统一错误响应；AppError 使用其自带状态码，其余为 500。
func Error(c *gin.Context, err error) {
	ae := AsAppError(err)
	if ae != nil {
		Fail(c, ae.HTTP, ae.Code, ae.Message)
		return
	}
	L(c).Error("unhandled error", "err", err)
	Fail(c, http.StatusInternalServerError, constants.CodeInternalError, constants.MessageOf(constants.CodeInternalError))
}

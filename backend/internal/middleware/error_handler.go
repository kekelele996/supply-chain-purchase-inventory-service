package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/util"
)

// ErrorHandler 统一异常处理中间件：捕获 panic，将错误转为统一响应。
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				util.L(c).Error(constants.LogTplPanicRecovered,
					"path", c.Request.URL.Path, "err", r, "stack", string(debug.Stack()))
				util.Fail(c, http.StatusInternalServerError, constants.CodeInternalError, constants.MessageOf(constants.CodeInternalError))
			}
		}()
		c.Next()
	}
}

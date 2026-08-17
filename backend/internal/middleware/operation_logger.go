package middleware

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/service"
	"github.com/supplychain/supplychain-api/internal/util"
)

var idPathPattern = regexp.MustCompile(`/(\d+)(/|$)`)

// OperationLogger 操作日志中间件：拦截 POST/PUT/DELETE 写操作，异步记录审计日志。
func OperationLogger(logService *service.OperationLogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Writer.Status() >= 400 {
			return
		}
		method := c.Request.Method
		var action constants.OperationAction
		switch method {
		case httpMethodPost:
			action = constants.ActionCreate
		case httpMethodPut:
			action = constants.ActionUpdate
		case httpMethodDelete:
			action = constants.ActionDelete
		default:
			return
		}
		path := c.Request.URL.Path
		if strings.Contains(path, "/auth/") {
			return
		}
		targetType := inferTargetType(path)
		if targetType == "" {
			return
		}
		targetID := inferTargetID(path)
		user := util.CurrentUser(c)
		detail := map[string]interface{}{
			"method": method,
			"path":   path,
			"status": c.Writer.Status(),
		}
		go logService.RecordOperation(context.Background(), user.ID, action, targetType, targetID, detail)
	}
}

func inferTargetType(path string) string {
	trimmed := strings.Trim(path, "/")
	parts := strings.Split(trimmed, "/")
	for _, p := range parts {
		switch p {
		case "suppliers", "inventory", "purchases", "users":
			return p
		}
	}
	return ""
}

func inferTargetID(path string) uint {
	m := idPathPattern.FindStringSubmatch(path)
	if len(m) == 2 {
		if id, err := strconv.ParseUint(m[1], 10, 64); err == nil {
			return uint(id)
		}
	}
	return 0
}

const (
	httpMethodPost   = "POST"
	httpMethodPut    = "PUT"
	httpMethodDelete = "DELETE"
)

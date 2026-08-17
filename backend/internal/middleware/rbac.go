package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/supplychain/supplychain-api/internal/constants"
	"github.com/supplychain/supplychain-api/internal/util"
)

// RequireRoles 角色权限中间件：当前用户角色必须在允许列表中。
func RequireRoles(roles ...constants.UserRole) gin.HandlerFunc {
	allowed := make(map[constants.UserRole]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		user := util.CurrentUser(c)
		if user.ID == 0 && user.Username == "" {
			util.Fail(c, http.StatusUnauthorized, constants.CodeUnauthorized, constants.MessageOf(constants.CodeUnauthorized))
			return
		}
		if !allowed[user.Role] {
			msg := fmt.Sprintf(constants.ErrTextForbidden, user.Role, c.Request.Method+" "+c.FullPath())
			util.Fail(c, http.StatusForbidden, constants.CodeForbidden, msg)
			return
		}
		c.Next()
	}
}

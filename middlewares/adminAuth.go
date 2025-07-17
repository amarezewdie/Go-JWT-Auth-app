package middlewares

import (
	"net/http"

	"go-jwt-mysql/utils"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware restricts access to admin users only
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get("userRole")
		if !exists {
			utils.RespondWithError(c, http.StatusForbidden, "Admin access required")
			c.Abort()
			return
		}

		role, ok := roleValue.(string)
		if !ok || role != "admin" {
			utils.RespondWithError(c, http.StatusForbidden, "Admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}

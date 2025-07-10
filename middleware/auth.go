package middleware

import (
	"net/http"
	"stock-dashboard/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(c *gin.Context) {
	if c.Request.Method == "OPTIONS" {
		c.Next()
		return
	}

	auth := c.GetHeader("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized: Missing or invalid token",
		})
		return
	}

	token := strings.TrimPrefix(auth, "Bearer ")
	claims, err := utils.VerifyToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized: Invalid token",
		})
		return
	}
	c.Set("userID", claims.UserID)
	c.Set("role", claims.Role)
	c.Next()

}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"message": "Forbidden: Admins only",
			})
			return
		}
		c.Next()
	}
}

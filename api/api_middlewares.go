package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func apiMiddlewareAuthSmtp() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		authorization = strings.TrimPrefix(authorization, "Bearer:")
		authorization = strings.TrimSpace(authorization)

		if authorization == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"error":   "Authorization header is required",
			})
			return
		}

		if authorization != config.API.Secret {
			c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]interface{}{
				"success": false,
				"error":   "Not authorized",
			})
			return
		}

		c.Next()
	}
}

func apiMiddlewareCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

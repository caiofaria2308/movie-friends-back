package api

import (
	"app/utils/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Lista de origens permitidas
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://localhost:8080",
			"http://0.0.0.0:3000",
			"http://0.0.0.0:5173",
			"http://0.0.0.0:8080",
			"https://movie-friends-front.onrender.com",
		}

		// Verifica se a origem está na lista de permitidas
		isAllowed := false
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin || strings.TrimSuffix(origin, "/") == allowedOrigin {
				isAllowed = true
				break
			}
		}

		if isAllowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := token.ValidateToken(token.ExtractToken(c))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RoleMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, err := token.ExtractTokenRole(c)
		if err != nil || role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}
		c.Next()
	}
}

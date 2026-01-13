package api

import (
	accounts_router "app/api/accounts"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(DB *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r = accounts_router.MountAccountsRouter(r, DB, AuthMiddleware())

	protected := r.Group("/api")
	protected.Use(AuthMiddleware())
	{

		admin := protected.Group("/admin")
		admin.Use(RoleMiddleware("admin"))
		{
			admin.GET("/dashboard", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "Welcome Admin"})
			})
		}
	}

	return r
}

package accounts_router

import (
	entity_accounts "app/entity/accounts"
	repository_accounts "app/infrascture/database/postgres/repository/accounts"
	usecase_accounts "app/usecase/accounts"
	"app/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=admin user guest"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type accountsRouter struct {
	usecase_user usecase_accounts.IUseCaseUser
}

func NewAccountsRouter(usecase_user usecase_accounts.IUseCaseUser) *accountsRouter {
	return &accountsRouter{usecase_user: usecase_user}
}

func (ar *accountsRouter) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := entity_accounts.User{
		Name:  input.Name,
		Email: input.Email,
		Role:  input.Role,
	}

	if err := user.EncryptedPassword(input.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	if err := ar.usecase_user.Register(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user_id": user.ID})
}

func (ar *accountsRouter) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := ar.usecase_user.Login(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	tokenString, err := token.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (ar *accountsRouter) GetMe(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := ar.usecase_user.FindById(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name": user.Name,
		"role": user.Role,
	})
}

func MountAccountsRouter(router *gin.Engine, DB *gorm.DB, authMiddleware gin.HandlerFunc) *gin.Engine {
	repo := repository_accounts.NewUserRepository(DB)
	usecase := usecase_accounts.NewUserUseCase(repo)
	// router group /auth
	ar := NewAccountsRouter(usecase)
	accounts := router.Group("/auth")
	{
		accounts.POST("/register", ar.Register)
		accounts.POST("/login", ar.Login)
	}
	// router group /api
	api := router.Group("/api")
	api.Use(authMiddleware)
	{
		api.GET("/user/profile", ar.GetMe)
	}
	return router
}

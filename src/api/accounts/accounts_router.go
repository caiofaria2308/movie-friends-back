package accounts_router

import (
	entity_accounts "app/entity/accounts"
	repository_accounts "app/infrascture/database/postgres/repository/accounts"
	usecase_accounts "app/usecase/accounts"
	"app/utils/token"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

type UserPixInput struct {
	PixKey string `json:"pix_key" binding:"required"`
}

type UserDayOffInput struct {
	InitHour    time.Time `json:"init_hour" binding:"required"`
	EndHour     time.Time `json:"end_hour" binding:"required"`
	Repeat      bool      `json:"repeat"`
	RepeatType  string    `json:"repeat_type"`
	RepeatValue string    `json:"repeat_value"`
}

type accountsRouter struct {
	usecase_user        usecase_accounts.IUseCaseUser
	usecase_user_pix    usecase_accounts.IUseCaseUserPix
	usecase_user_dayoff usecase_accounts.IUseCaseUserDayOff
}

func NewAccountsRouter(usecase_user usecase_accounts.IUseCaseUser, usecase_user_pix usecase_accounts.IUseCaseUserPix, usecase_user_dayoff usecase_accounts.IUseCaseUserDayOff) *accountsRouter {
	return &accountsRouter{
		usecase_user:        usecase_user,
		usecase_user_pix:    usecase_user_pix,
		usecase_user_dayoff: usecase_user_dayoff,
	}
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

func (ar *accountsRouter) CreatePix(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input UserPixInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userPix := entity_accounts.UserPix{
		PixKey: input.PixKey,
	}

	if err := ar.usecase_user_pix.Create(&userPix, userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, userPix)
}

func (ar *accountsRouter) GetPix(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	userPix, err := ar.usecase_user_pix.GetById(id, userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userPix)
}

func (ar *accountsRouter) ListPix(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userPixs, err := ar.usecase_user_pix.GetAll(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userPixs)
}

func (ar *accountsRouter) DeletePix(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	if err := ar.usecase_user_pix.Delete(id, userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pix key deleted successfully"})
}

func (ar *accountsRouter) CreateDayOff(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input UserDayOffInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dayOff := entity_accounts.UserDayOff{
		InitHour:    &input.InitHour,
		EndHour:     &input.EndHour,
		Repeat:      input.Repeat,
		RepeatType:  input.RepeatType,
		RepeatValue: input.RepeatValue,
	}

	if err := ar.usecase_user_dayoff.Create(&dayOff, userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dayOff)
}

func (ar *accountsRouter) ListDayOff(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse filter query parameters
	filterType := c.Query("filter_type")
	yearStr := c.Query("year")
	weekStr := c.Query("week")
	monthStr := c.Query("month")

	var year, week, month int

	// Parse year if provided
	if yearStr != "" {
		year, err = strconv.Atoi(yearStr)
		if err != nil || year < 1900 || year > 3000 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year parameter"})
			return
		}
	}

	// Parse week if provided
	if weekStr != "" {
		week, err = strconv.Atoi(weekStr)
		if err != nil || week < 1 || week > 53 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid week parameter (must be 1-53)"})
			return
		}
	}

	// Parse month if provided
	if monthStr != "" {
		month, err = strconv.Atoi(monthStr)
		if err != nil || month < 1 || month > 12 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month parameter (must be 1-12)"})
			return
		}
	}

	// Validate filter combinations
	if filterType != "" {
		switch filterType {
		case "week":
			if year == 0 || week == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Week filter requires both 'year' and 'week' parameters"})
				return
			}
		case "month":
			if year == 0 || month == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Month filter requires both 'year' and 'month' parameters"})
				return
			}
		case "year":
			if year == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Year filter requires 'year' parameter"})
				return
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter_type. Must be 'week', 'month', or 'year'"})
			return
		}
	}

	dayOffs, err := ar.usecase_user_dayoff.GetAll(userId, filterType, year, week, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dayOffs)
}

func (ar *accountsRouter) UpdateDayOff(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	// Mode param: all, single, future
	mode := c.DefaultQuery("mode", "single")

	var input UserDayOffInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dayOff := entity_accounts.UserDayOff{
		ID:       &id,
		InitHour: &input.InitHour,
		EndHour:  &input.EndHour,
	}

	if err := ar.usecase_user_dayoff.Update(&dayOff, userId, mode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Day off updated successfully"})
}

func (ar *accountsRouter) DeleteDayOff(c *gin.Context) {
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	// Mode param: all, single, future
	mode := c.DefaultQuery("mode", "single")

	if err := ar.usecase_user_dayoff.Delete(id, userId, mode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Day off deleted successfully"})
}

func MountAccountsRouter(router *gin.Engine, DB *gorm.DB, authMiddleware gin.HandlerFunc) *gin.Engine {
	repoUser := repository_accounts.NewUserRepository(DB)
	usecaseUser := usecase_accounts.NewUserUseCase(repoUser)

	repoPix := repository_accounts.NewUserPixRepository(DB)
	usecasePix := usecase_accounts.NewUserPixUseCase(repoPix)

	repoDayOff := repository_accounts.NewUserDayOffRepository(DB)
	usecaseDayOff := usecase_accounts.NewUserDayOffUseCase(repoDayOff)

	// router group /auth
	ar := NewAccountsRouter(usecaseUser, usecasePix, usecaseDayOff)
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

		// Pix Routes
		api.POST("/user/pix", ar.CreatePix)
		api.GET("/user/pix", ar.ListPix)
		api.GET("/user/pix/:id", ar.GetPix)
		api.DELETE("/user/pix/:id", ar.DeletePix)

		// DayOff Routes
		api.POST("/user/dayoff", ar.CreateDayOff)
		api.GET("/user/dayoff", ar.ListDayOff)
		api.PUT("/user/dayoff/:id", ar.UpdateDayOff)
		api.DELETE("/user/dayoff/:id", ar.DeleteDayOff)
	}
	return router
}

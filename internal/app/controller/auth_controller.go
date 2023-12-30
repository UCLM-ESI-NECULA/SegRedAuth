package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"seg-red-auth/internal/app/common"
	"seg-red-auth/internal/app/dao"
	"seg-red-auth/internal/app/repository"
	"seg-red-auth/internal/app/service"
)

type AuthControllerImpl struct {
	svc service.AuthService
}

func NewAuthController() *AuthControllerImpl {
	userRepo := repository.NewUserRepository(common.ConnectToDB())
	return &AuthControllerImpl{svc: service.NewAuthService(userRepo)}
}

type AuthController interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
}

// RegisterRoutes registers the authentication routes
func (ac *AuthControllerImpl) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/signup", ac.Signup)
	router.POST("/login", ac.Login)
}

// Signup handles the /signup endpoint
func (ac *AuthControllerImpl) Signup(c *gin.Context) {
	user, err := checkUserInput(c)
	if err != nil {
		common.NewAPIError(c, http.StatusBadRequest, err, err.Error())
		return
	}
	token := ac.svc.CreateUser(user.Username, user.Password)
	c.JSON(http.StatusOK, gin.H{"access_token": token})
}

// Login handles the /login endpoint
func (ac *AuthControllerImpl) Login(c *gin.Context) {
	user, err := checkUserInput(c)
	if err != nil {
		common.NewAPIError(c, http.StatusBadRequest, err, err.Error())
		return
	}
	token, err := ac.svc.AuthenticateUser(user.Username, user.Password)
	if err != nil {
		common.NewAPIError(c, http.StatusUnauthorized, err, "invalid credentials")
		return
	}
	c.JSON(http.StatusOK, gin.H{"access_token": token})
}

// checkUserInput checks if the user input is valid
func checkUserInput(c *gin.Context) (dao.User, error) {
	var user dao.User
	if err := c.ShouldBindJSON(&user); err != nil {
		common.NewAPIError(c, http.StatusUnauthorized, err, "error when mapping request")
		return user, err
	}
	if user.Username == "" || user.Password == "" {
		return user, fmt.Errorf("username and password are required")
	}
	return user, nil
}

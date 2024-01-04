package controller

import (
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

func NewAuthController(g *gin.RouterGroup) *AuthControllerImpl {
	userRepo, err := repository.NewUserRepository(common.ConnectToDB())
	if err != nil {
		panic(err)
	}
	c := &AuthControllerImpl{svc: service.NewAuthService(userRepo)}
	c.RegisterRoutes(g)
	return c
}

type AuthController interface {
	Signup(c *gin.Context)
	Login(c *gin.Context)
	CheckToken(c *gin.Context)
}

// RegisterRoutes registers the authentication routes
func (ac *AuthControllerImpl) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/signup", ac.Signup)
	router.POST("/login", ac.Login)
	router.POST("/checkToken", ac.CheckToken)
}

// Signup handles the /signup endpoint
func (ac *AuthControllerImpl) Signup(c *gin.Context) {
	// Check input
	user, err := checkUserInput(c)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	// Create user
	token, err := ac.svc.CreateUser(user.Username, user.Password)
	if err != nil {
		common.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, dao.Token{Token: token})
}

// Login handles the /login endpoint
func (ac *AuthControllerImpl) Login(c *gin.Context) {
	// Check input
	user, err := checkUserInput(c)
	if err != nil {
		common.HandleError(c, err)
		return
	}

	// Login user
	token, err := ac.svc.AuthenticateUser(user.Username, user.Password)
	if err != nil {
		common.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, dao.Token{Token: token})
}

func (ac *AuthControllerImpl) CheckToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		common.ForwardError(c, common.UnauthorizedError("authorization header is required"))
		return
	}

	// Login user
	username, err := ac.svc.ValidateToken(token)
	if err != nil {
		common.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, dao.User{Username: username})
}

// checkUserInput checks if the user input is valid
func checkUserInput(c *gin.Context) (*dao.User, error) {
	var user *dao.User
	if err := c.ShouldBindJSON(&user); err != nil {
		return nil, err
	}
	if user.Username == "" || user.Password == "" {
		return nil, common.EmptyParamsError("username")
	}
	if user.Password == "" {
		return nil, common.EmptyParamsError("password")
	}
	return user, nil
}

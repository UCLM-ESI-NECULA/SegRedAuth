package service

import (
	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"os"
	"seg-red-auth/internal/app/auth"
	"seg-red-auth/internal/app/common"
	"seg-red-auth/internal/app/dao"
	"seg-red-auth/internal/app/repository"
	"time"
)

type AuthServiceImpl struct {
	ur repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *AuthServiceImpl {
	return &AuthServiceImpl{
		ur: repo,
	}
}

type AuthService interface {
	CreateUser(username, password string) (string, error)
	AuthenticateUser(username, password string) (string, error)
	GetVersion() map[string]string
	ValidateToken(tokenString string) (string, error)
}

func (svc *AuthServiceImpl) GetVersion() map[string]string {
	return map[string]string{
		"version": "1.0",
	}
}

func (svc *AuthServiceImpl) CreateUser(username, password string) (string, error) {
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return "", err
	}

	user := &dao.User{Username: username, Password: hashedPassword} // create a pointer
	save, err := svc.ur.Save(user)
	if err != nil {
		return "", err
	}

	// Generate token
	token, err := auth.GenerateToken(save.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (svc *AuthServiceImpl) AuthenticateUser(username, password string) (string, error) {
	user, userError := svc.ur.GetUser(username)
	if userError != nil {
		return "", userError
	}

	err := auth.CheckPasswordHash(password, user.Password)
	if err != nil {
		return "", common.UnauthorizedError("invalid password")
	}

	// Generate token
	token, err := auth.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

// ValidateToken checks if the provided token string is valid and returns the corresponding user.
func (svc *AuthServiceImpl) ValidateToken(tokenString string) (string, error) {
	secret := os.Getenv("JWT_SECRET_KEY")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Error("Unexpected signing method: %v", token.Header["alg"])
			return nil, common.UnauthorizedError("Unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", common.UnauthorizedError(err.Error())
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok := claims["username"].(string)
		if !ok {
			log.Error("username not found in token")
			return "", common.UnauthorizedError("username not found in token")
		}

		if exp, ok := claims["exp"].(float64); ok {
			if time.Unix(int64(exp), 0).Before(time.Now()) {
				return "", common.UnauthorizedError("token expired")
			}
		} else {
			return "", common.UnauthorizedError("exp field not found in token")
		}

		return username, nil
	}

	return "", common.UnauthorizedError("invalid token")
}

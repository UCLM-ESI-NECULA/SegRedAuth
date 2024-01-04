package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"seg-red-auth/internal/app/dao"
)

func NewUserRepository(db *gorm.DB) (*UserRepositoryImpl, error) {
	err := db.AutoMigrate(&dao.User{})
	if err != nil {
		return nil, fmt.Errorf("error when migrating user db: %v", err)
	}
	return &UserRepositoryImpl{db: db}, nil
}

type UserRepository interface {
	FindAllUser() (*[]dao.User, error)
	Save(user *dao.User) (*dao.User, error)
	DeleteUserById(id int) error
	GetUser(username string) (*dao.User, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func (ur UserRepositoryImpl) FindAllUser() (*[]dao.User, error) {
	var users *[]dao.User
	result := ur.db.Find(&users)
	if result.Error != nil {
		_ = fmt.Errorf("error when get all user: %v", result.Error)
	}
	return users, nil
}

func (ur UserRepositoryImpl) GetUser(username string) (*dao.User, error) {
	var user *dao.User
	result := ur.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error when getting user: %v", result.Error)
	}
	return user, nil
}
func (ur UserRepositoryImpl) Save(user *dao.User) (*dao.User, error) {
	result := ur.db.Save(user)
	if result.Error != nil {
		return nil, fmt.Errorf("error when saving user: %v", result.Error)
	}
	return user, nil
}

func (ur UserRepositoryImpl) DeleteUserById(id int) error {
	result := ur.db.Delete(&dao.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("error when deleting user: %v", result.Error)
	}
	return nil
}

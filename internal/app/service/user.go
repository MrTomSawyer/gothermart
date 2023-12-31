package service

import (
	"github.com/MrTomSawyer/loyalty-system/internal/app/apperrors/autherr"
	"github.com/MrTomSawyer/loyalty-system/internal/app/config"
	"github.com/MrTomSawyer/loyalty-system/internal/app/entity"
	"github.com/MrTomSawyer/loyalty-system/internal/app/interfaces"
	"github.com/MrTomSawyer/loyalty-system/internal/app/logger"
	"github.com/MrTomSawyer/loyalty-system/internal/app/models"
)

type UserService struct {
	cfg            *config.Config
	UserRepository interfaces.UserRepository
	authService    interfaces.Auth
}

func NewUserService(cfg *config.Config, userRepo interfaces.UserRepository, authService interfaces.Auth) *UserService {
	return &UserService{
		cfg:            cfg,
		UserRepository: userRepo,
		authService:    authService,
	}
}

func (u *UserService) Login(user models.User) (string, error) {
	usr, err := u.UserRepository.GetUserByLogin(user.Login)
	if err != nil {
		return "", err
	}
	userEntity := entity.User{
		ID:    usr.ID,
		Login: usr.Login,
	}
	userEntity.SetRowPasswordString(usr.Password)

	isCorrect := userEntity.ComparePassword(user.Password)
	if !isCorrect {
		logger.Log.Error("wrong password")
		return "", autherr.ErrWrongCredentials
	}

	jwt, err := u.authService.JWT(userEntity.ID)
	if err != nil {
		logger.Log.Errorf("failed to sign jwt")
		return "", autherr.ErrJWTSignFailure
	}

	return jwt, nil
}

func (u *UserService) CreateUser(user models.User) (string, error) {
	newUser := entity.User{
		Login:     user.Login,
		Balance:   0,
		Withdrawn: 0,
	}
	newUser.HashPassword(user.Password)

	id, err := u.UserRepository.CreateUser(newUser)
	if err != nil {
		return "", err
	}

	jwt, err := u.authService.JWT(id)
	if err != nil {
		logger.Log.Errorf("failed to sign jwt")
		return "", autherr.ErrJWTSignFailure
	}

	return jwt, nil
}

func (u *UserService) GetUserBalance(userID int) (models.Balance, error) {
	return u.UserRepository.GetUserBalance(userID)
}

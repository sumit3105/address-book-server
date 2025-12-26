package service

import (
	appError "address-book-server/error"
	"address-book-server/logger"
	"address-book-server/model"
	"address-book-server/repository"
	"address-book-server/utils"
	"errors"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(email, password string) error
	Login(email, password string) (string, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo}
}

func (service *authService) Register(email, password string) error {

	logger.Log.Info(
		"Registrating user started",
		zap.String("email", email),
	)

	_, err := service.repo.FindByEmail(email)
	
	if err == nil {
		
		logger.Log.Error(
			"User exists",
			zap.String("email", email),

		)

		return appError.BadRequest(
			"User already exists",
			nil,
		)
	}
	
	if !errors.Is(err, gorm.ErrRecordNotFound) {

		logger.Log.Error(
			"Internal server error",
			zap.String("error", err.Error()),

		)

		return appError.Internal(
			"Internal server error",
			err,
		)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {

		logger.Log.Error(
			"Error generating password hash",
			zap.String("error", err.Error()),

		)

		return appError.BadRequest(
			"Error generating password hash",
			err,
		)
	}

	user := model.User{
		Email: email,
		PasswordHash: string(hash),
	}

	logger.Log.Info(
		"New User registered",
		zap.Uint64("id", user.ID),
		zap.String("email", user.Email),
	)

	return service.repo.Create(&user)
}

func (service *authService) Login(email, password string) (string, error) {
	
	logger.Log.Info(
		"User trying to login...",
		zap.String("email", email),
	)

	user, err := service.repo.FindByEmail(email)
	if err != nil {

		logger.Log.Error(
			"Invalid email",
			zap.String("error", err.Error()),

		)

		return "", appError.Forbidden(
			"Invalid Email",
			err,
		)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {

		logger.Log.Error(
			"Invalid password",
			zap.String("error", err.Error()),

		)

		return "", appError.Forbidden(
			"Invalid Password",
			err,
		)
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {

		logger.Log.Error(
			"Error in generating token",
			zap.String("error", err.Error()),

		)

		return "", appError.Internal(
			"Error generating token",
			err,
		)
	}

	logger.Log.Info(
		"User logged in",
		zap.String("email", email),
	)

	return token, nil

}
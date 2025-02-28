package auth

import (
	"errors"
	"url-shortener/internal/user"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	*user.UserRepository
}

func NewAuthService(userRepo user.UserRepository) *AuthService {
	return &AuthService{
		UserRepository: &userRepo,
	}
}

func (svc *AuthService) Login(email string, password string) (string, error) {
	findUser, err := svc.UserRepository.FindByEmail(email)

	if findUser == nil {
		return "", errors.New(ErrUserNotFound)
	}

	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(findUser.Password), []byte(password))

	if err != nil {
		return "", errors.New(ErrWrongCredentials)
	}

	return findUser.Email, nil
}

func (svc *AuthService) Register(name string, email string, password string) (string, error) {
	existedUser, _ := svc.UserRepository.FindByEmail(email)

	if existedUser != nil {
		return "", errors.New(ErrUserExists)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	user, err := svc.UserRepository.Create(&user.User{
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	})

	if err != nil {
		return "", err
	}

	return user.Email, nil
}

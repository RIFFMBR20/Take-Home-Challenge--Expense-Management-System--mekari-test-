package services

import (
	"backend-test-mekari/internal/repositories"
	"errors"
)

type AuthService interface {
	Login(email, password string) (string, string, uint, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(ur repositories.UserRepository) AuthService {
	return &authService{userRepo: ur}
}

func (s *authService) Login(email, password string) (string, string, uint, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", "", 0, errors.New("user not found")
	}

	if user.Password != password {
		return "", "", 0, errors.New("wrong password")
	}

	token := "mock-jwt-token-for-" + user.Email

	return token, user.Role, user.ID, nil
}

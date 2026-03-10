package service

import (
	"bambu-farm/domain"
	"bambu-farm/pkg/auth"
	"bambu-farm/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepo *repository.AuthRepository
}

func NewAuthService(authRepo *repository.AuthRepository) *AuthService {
	return &AuthService{authRepo: authRepo}
}

func (s *AuthService) RegisterUser(email, password, orgName string) (*domain.User, error) {
	existingUser, err := s.authRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	org, err := s.authRepo.GetOrganizationByName(orgName)
	if err != nil {
		return nil, err
	}
	
	if org == nil {
		org = &domain.Organization{Name: orgName}
		if err := s.authRepo.CreateOrganization(org); err != nil {
			return nil, err
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &domain.User{
		Email:          email,
		PasswordHash:   string(hashedPassword),
		OrganizationID: org.ID,
	}

	if err := s.authRepo.CreateUser(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *AuthService) LoginUser(email, password string) (string, string, error) {
	user, err := s.authRepo.GetUserByEmail(email)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	var roles []string
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}

	accessToken, refreshToken, err := auth.GenerateTokenPair(user.ID, user.OrganizationID, user.Email, roles)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) GetUser(email string) (*domain.User, error) {
	return s.authRepo.GetUserByEmail(email)
}

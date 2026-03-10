package repository

import (
	"bambu-farm/domain"
	"errors"

	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	result := r.db.Preload("Organization").Preload("Roles.Permissions").Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *AuthRepository) CreateUser(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *AuthRepository) GetOrganizationByName(name string) (*domain.Organization, error) {
	var org domain.Organization
	result := r.db.Where("name = ?", name).First(&org)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &org, nil
}

func (r *AuthRepository) CreateOrganization(org *domain.Organization) error {
	return r.db.Create(org).Error
}

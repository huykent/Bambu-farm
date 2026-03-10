package repository

import (
	"bambu-farm/domain"
	"errors"

	"gorm.io/gorm"
)

type PrinterRepository struct {
	db *gorm.DB
}

func NewPrinterRepository(db *gorm.DB) *PrinterRepository {
	return &PrinterRepository{db: db}
}

func (r *PrinterRepository) Create(printer *domain.Printer) error {
	return r.db.Create(printer).Error
}

func (r *PrinterRepository) FindAll(orgID uint) ([]domain.Printer, error) {
	var printers []domain.Printer
	result := r.db.Where("organization_id = ?", orgID).Find(&printers)
	return printers, result.Error
}

func (r *PrinterRepository) FindByID(id uint, orgID uint) (*domain.Printer, error) {
	var printer domain.Printer
	result := r.db.Where("id = ? AND organization_id = ?", id, orgID).First(&printer)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &printer, nil
}

func (r *PrinterRepository) Delete(id uint, orgID uint) error {
	result := r.db.Where("id = ? AND organization_id = ?", id, orgID).Delete(&domain.Printer{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("printer not found")
	}
	return nil
}

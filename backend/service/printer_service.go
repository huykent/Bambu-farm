package service

import (
	"bambu-farm/domain"
	"bambu-farm/repository"
	"errors"
)

type PrinterService struct {
	printerRepo *repository.PrinterRepository
}

func NewPrinterService(printerRepo *repository.PrinterRepository) *PrinterService {
	return &PrinterService{printerRepo: printerRepo}
}

func (s *PrinterService) AddPrinter(orgID uint, printerID, name, ipAddress, accessToken, model string) (*domain.Printer, error) {
	// Simple validation, would be expanded
	if printerID == "" || name == "" || ipAddress == "" {
		return nil, errors.New("missing required printer fields")
	}

	printer := &domain.Printer{
		OrganizationID: orgID,
		PrinterID:      printerID,
		Name:           name,
		IPAddress:      ipAddress,
		AccessToken:    accessToken, // In a real app, encrypt this before saving
		Model:          model,
		Status:         "offline", // Initial status
	}

	if err := s.printerRepo.Create(printer); err != nil {
		return nil, err
	}

	return printer, nil
}

func (s *PrinterService) ListPrinters(orgID uint) ([]domain.Printer, error) {
	return s.printerRepo.FindAll(orgID)
}

func (s *PrinterService) GetPrinter(id uint, orgID uint) (*domain.Printer, error) {
	printer, err := s.printerRepo.FindByID(id, orgID)
	if err != nil {
		return nil, err
	}
	if printer == nil {
		return nil, errors.New("printer not found")
	}
	return printer, nil
}

func (s *PrinterService) DeletePrinter(id uint, orgID uint) error {
	return s.printerRepo.Delete(id, orgID)
}

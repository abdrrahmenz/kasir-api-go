package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetTodayReport() (*models.DailySalesReport, error) {
	return s.repo.GetTodayReport()
}

func (s *ReportService) GetDateRangeReport(startDate, endDate string) (*models.DateRangeReport, error) {
	return s.repo.GetDateRangeReport(startDate, endDate)
}

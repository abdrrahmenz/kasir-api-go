package repositories

import (
	"database/sql"
	"kasir-api/models"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) GetTodayReport() (*models.DailySalesReport, error) {
	report := &models.DailySalesReport{}

	err := repo.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0), COUNT(*)
		FROM transactions
		WHERE DATE(created_at) = CURRENT_DATE
	`).Scan(&report.TotalRevenue, &report.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	bestSelling := &models.BestSellingProduct{}
	err = repo.db.QueryRow(`
		SELECT p.name, COALESCE(SUM(td.quantity), 0) as qty
		FROM transaction_details td
		JOIN products p ON p.id = td.product_id
		JOIN transactions t ON t.id = td.transaction_id
		WHERE DATE(t.created_at) = CURRENT_DATE
		GROUP BY p.name
		ORDER BY qty DESC
		LIMIT 1
	`).Scan(&bestSelling.Nama, &bestSelling.QtyTerjual)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == nil {
		report.ProdukTerlaris = bestSelling
	}

	return report, nil
}

func (repo *ReportRepository) GetDateRangeReport(startDate, endDate string) (*models.DateRangeReport, error) {
	report := &models.DateRangeReport{
		StartDate: startDate,
		EndDate:   endDate,
	}

	err := repo.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount), 0), COUNT(*)
		FROM transactions
		WHERE DATE(created_at) BETWEEN $1 AND $2
	`, startDate, endDate).Scan(&report.TotalRevenue, &report.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	bestSelling := &models.BestSellingProduct{}
	err = repo.db.QueryRow(`
		SELECT p.name, COALESCE(SUM(td.quantity), 0) as qty
		FROM transaction_details td
		JOIN products p ON p.id = td.product_id
		JOIN transactions t ON t.id = td.transaction_id
		WHERE DATE(t.created_at) BETWEEN $1 AND $2
		GROUP BY p.name
		ORDER BY qty DESC
		LIMIT 1
	`, startDate, endDate).Scan(&bestSelling.Nama, &bestSelling.QtyTerjual)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == nil {
		report.ProdukTerlaris = bestSelling
	}

	return report, nil
}

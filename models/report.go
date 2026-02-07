package models

type BestSellingProduct struct {
Nama       string `json:"nama"`
QtyTerjual int    `json:"qty_terjual"`
}

type DailySalesReport struct {
TotalRevenue   int                 `json:"total_revenue"`
TotalTransaksi int                 `json:"total_transaksi"`
ProdukTerlaris *BestSellingProduct `json:"produk_terlaris"`
}

type DateRangeReport struct {
StartDate      string              `json:"start_date"`
EndDate        string              `json:"end_date"`
TotalRevenue   int                 `json:"total_revenue"`
TotalTransaksi int                 `json:"total_transaksi"`
ProdukTerlaris *BestSellingProduct `json:"produk_terlaris"`
}
